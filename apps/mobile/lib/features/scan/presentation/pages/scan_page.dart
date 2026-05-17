import 'dart:math' as math;

import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../app/theme/app_theme.dart';
import '../controllers/scan_controller.dart';
import '../widgets/scanner_preview_frame.dart';

class ScanPage extends ConsumerStatefulWidget {
  const ScanPage({super.key});

  @override
  ConsumerState<ScanPage> createState() => _ScanPageState();
}

class _ScanPageState extends ConsumerState<ScanPage> {
  CameraController? _cameraController;
  String? _cameraError;

  @override
  void initState() {
    super.initState();
    _initializeCamera();
  }

  @override
  void dispose() {
    _cameraController?.dispose();
    super.dispose();
  }

  Future<void> _initializeCamera() async {
    setState(() {
      _cameraError = null;
    });

    try {
      final cameras = await availableCameras();
      if (cameras.isEmpty) {
        if (!mounted) return;
        setState(() {
          _cameraError = 'Tidak ada kamera yang terdeteksi di perangkat ini.';
        });
        return;
      }

      final selectedCamera = cameras.firstWhere(
        (camera) => camera.lensDirection == CameraLensDirection.back,
        orElse: () => cameras.first,
      );

      final controller = CameraController(
        selectedCamera,
        ResolutionPreset.high,
        enableAudio: false,
        imageFormatGroup: ImageFormatGroup.jpeg,
      );

      await _cameraController?.dispose();
      _cameraController = controller;
      await controller.initialize();

      if (!mounted) return;
      setState(() {});
    } on CameraException catch (error) {
      if (!mounted) return;
      setState(() {
        _cameraError = _cameraErrorMessage(error);
      });
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _cameraError = 'Kamera tidak dapat dibuka saat ini.';
      });
    }
  }

  String _cameraErrorMessage(CameraException error) {
    return switch (error.code) {
      'CameraAccessDenied' =>
        'Izin kamera ditolak. Aktifkan izin kamera untuk NutriScan.',
      'CameraAccessDeniedWithoutPrompt' =>
        'Izin kamera belum aktif. Buka pengaturan perangkat untuk mengaktifkannya.',
      'CameraAccessRestricted' =>
        'Akses kamera dibatasi oleh pengaturan perangkat.',
      _ => error.description ?? 'Kamera tidak dapat dibuka saat ini.',
    };
  }

  Future<void> _captureAndAnalyze() async {
    final controller = _cameraController;
    if (controller == null || !controller.value.isInitialized) {
      ref
          .read(scanControllerProvider.notifier)
          .fail('Kamera belum siap. Coba lagi sebentar.');
      return;
    }

    if (controller.value.isTakingPicture) return;

    try {
      final file = await controller.takePicture();
      await ref
          .read(scanControllerProvider.notifier)
          .analyzeImage(capturedImagePath: file.path);
    } on CameraException catch (error) {
      ref
          .read(scanControllerProvider.notifier)
          .fail(error.description ?? 'Foto tidak berhasil diambil.');
    }
  }

  @override
  Widget build(BuildContext context) {
    final scanState = ref.watch(scanControllerProvider);

    return Scaffold(
      backgroundColor: Colors.black,
      body: Stack(
        children: [
          Positioned.fill(
            child: ScannerPreviewFrame(
              controller: _cameraController,
              errorMessage: _cameraError,
              onRetry: _initializeCamera,
            ),
          ),
          const Positioned.fill(child: _CameraReadabilityOverlay()),
          SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(20, 12, 20, 24),
              child: Column(
                children: [
                  _ScanTopBar(onClose: () => context.pop()),
                  const SizedBox(height: 40),
                  Expanded(child: _ScanGuide(scanState: scanState)),
                  _ScanBottomAction(
                    scanState: scanState,
                    isCameraReady:
                        _cameraController?.value.isInitialized == true &&
                        _cameraError == null,
                    onCapture: _captureAndAnalyze,
                    onReset: () => ref
                        .read(scanControllerProvider.notifier)
                        .reset(),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _CameraReadabilityOverlay extends StatelessWidget {
  const _CameraReadabilityOverlay();

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Colors.black.withValues(alpha: 0.38),
            Colors.transparent,
            Colors.black.withValues(alpha: 0.5),
          ],
          stops: const [0, 0.46, 1],
        ),
      ),
    );
  }
}

class _ScanTopBar extends StatelessWidget {
  final VoidCallback onClose;

  const _ScanTopBar({required this.onClose});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        _RoundIconButton(
          icon: Icons.close,
          tooltip: 'Tutup scanner',
          onPressed: onClose,
        ),
        const Spacer(),
        Text(
          'aura plate',
          style: Theme.of(context).textTheme.titleMedium?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.w700,
          ),
        ),
        const Spacer(),
        const SizedBox(width: 48),
      ],
    );
  }
}

class _ScanGuide extends StatelessWidget {
  final ScanState scanState;

  const _ScanGuide({required this.scanState});

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final frameSize = math.min(constraints.maxWidth - 16, 292.0);

        return Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Stack(
              clipBehavior: Clip.none,
              children: [
                SizedBox.square(
                  dimension: frameSize,
                  child: CustomPaint(
                    painter: _DashedFramePainter(color: AppColors.darkNavy),
                  ),
                ),
                if (scanState.status == ScanStatus.success) ...[
                  const Positioned(
                    left: 26,
                    top: 62,
                    child: _NutritionBadge(
                      label: 'Vitamin',
                      value: '30 kcal',
                    ),
                  ),
                  const Positioned(
                    right: 18,
                    top: 44,
                    child: _NutritionBadge(
                      label: 'Protein Utama',
                      value: '250 kcal',
                    ),
                  ),
                  const Positioned(
                    right: -4,
                    top: 114,
                    child: _NutritionBadge(
                      label: 'Karbohidrat',
                      value: '200 kcal',
                    ),
                  ),
                  const Positioned(
                    right: 30,
                    bottom: 62,
                    child: _NutritionBadge(
                      label: 'Serat & Mineral',
                      value: '100 kcal',
                    ),
                  ),
                ],
              ],
            ),
            AnimatedSwitcher(
              duration: const Duration(milliseconds: 220),
              child: scanState.status == ScanStatus.success
                  ? const _EstimateSummary()
                  : const SizedBox(height: 84),
            ),
          ],
        );
      },
    );
  }
}

class _ScanBottomAction extends StatelessWidget {
  final ScanState scanState;
  final bool isCameraReady;
  final VoidCallback onCapture;
  final VoidCallback onReset;

  const _ScanBottomAction({
    required this.scanState,
    required this.isCameraReady,
    required this.onCapture,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    if (scanState.status == ScanStatus.analyzing) {
      return const _AnalyzingPill();
    }

    if (scanState.status == ScanStatus.success) {
      return _RoundIconButton(
        icon: Icons.refresh,
        tooltip: 'Scan ulang',
        size: 58,
        onPressed: onReset,
      );
    }

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        if (scanState.status == ScanStatus.error) ...[
          _InlineError(message: scanState.errorMessage ?? 'Scan gagal.'),
          const SizedBox(height: 14),
        ],
        _RoundIconButton(
          icon: Icons.camera_alt,
          tooltip: 'Ambil foto makanan',
          size: 58,
          onPressed: isCameraReady ? onCapture : null,
        ),
      ],
    );
  }
}

class _RoundIconButton extends StatelessWidget {
  final IconData icon;
  final String tooltip;
  final VoidCallback? onPressed;
  final double size;

  const _RoundIconButton({
    required this.icon,
    required this.tooltip,
    required this.onPressed,
    this.size = 48,
  });

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: tooltip,
      child: Material(
        color: onPressed == null
            ? Colors.white.withValues(alpha: 0.26)
            : AppColors.darkNavy,
        shape: const CircleBorder(),
        clipBehavior: Clip.antiAlias,
        child: InkWell(
          onTap: onPressed,
          child: SizedBox.square(
            dimension: size,
            child: Icon(
              icon,
              color: Colors.white,
              size: size * 0.46,
            ),
          ),
        ),
      ),
    );
  }
}

class _AnalyzingPill extends StatelessWidget {
  const _AnalyzingPill();

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.62),
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: Colors.white.withValues(alpha: 0.18)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 12),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const SizedBox.square(
              dimension: 18,
              child: CircularProgressIndicator(
                strokeWidth: 2.4,
                color: AppColors.primaryGreen,
              ),
            ),
            const SizedBox(width: 12),
            Text(
              'Menganalisis nutrisi...',
              style: Theme.of(context).textTheme.labelLarge?.copyWith(
                color: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _EstimateSummary extends StatelessWidget {
  const _EstimateSummary();

  @override
  Widget build(BuildContext context) {
    return Padding(
      key: const ValueKey('estimate-summary'),
      padding: const EdgeInsets.only(top: 12),
      child: Column(
        children: [
          Text(
            'Total Estimasi:',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: AppColors.primaryGreen,
              fontWeight: FontWeight.w800,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            '~580 kcal',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
              color: Colors.white,
              fontWeight: FontWeight.w800,
            ),
          ),
          const SizedBox(height: 8),
          DecoratedBox(
            decoration: BoxDecoration(
              color: AppColors.primaryGreen,
              borderRadius: BorderRadius.circular(999),
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              child: Text(
                'Show Aura Plate',
                style: Theme.of(context).textTheme.labelSmall?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _InlineError extends StatelessWidget {
  final String message;

  const _InlineError({required this.message});

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.62),
        borderRadius: BorderRadius.circular(18),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
        child: Text(
          message,
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
}

class _NutritionBadge extends StatelessWidget {
  final String label;
  final String value;

  const _NutritionBadge({
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        color: AppColors.primaryGreen,
        borderRadius: BorderRadius.circular(4),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.18),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 4),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              '$label:',
              style: const TextStyle(
                color: AppColors.darkNavy,
                fontSize: 8,
                fontWeight: FontWeight.w700,
                height: 1,
              ),
            ),
            const SizedBox(height: 1),
            Text(
              value,
              style: const TextStyle(
                color: AppColors.darkNavy,
                fontSize: 10,
                fontWeight: FontWeight.w800,
                height: 1,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _DashedFramePainter extends CustomPainter {
  final Color color;

  const _DashedFramePainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 3
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round;

    final rect = Rect.fromLTWH(2, 2, size.width - 4, size.height - 4);
    final path = Path()
      ..addRRect(RRect.fromRectAndRadius(rect, const Radius.circular(10)));

    for (final metric in path.computeMetrics()) {
      var distance = 0.0;
      while (distance < metric.length) {
        final end = math.min(distance + 18, metric.length);
        canvas.drawPath(metric.extractPath(distance, end), paint);
        distance += 34;
      }
    }
  }

  @override
  bool shouldRepaint(covariant _DashedFramePainter oldDelegate) {
    return oldDelegate.color != color;
  }
}
