import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

import '../../../../app/theme/app_theme.dart';

class ScannerPreviewFrame extends StatelessWidget {
  final CameraController? controller;
  final String? errorMessage;
  final VoidCallback? onRetry;

  const ScannerPreviewFrame({
    super.key,
    required this.controller,
    this.errorMessage,
    this.onRetry,
  });

  @override
  Widget build(BuildContext context) {
    final cameraController = controller;

    if (errorMessage != null) {
      return _CameraMessage(
        icon: Icons.no_photography_outlined,
        title: 'Kamera belum bisa dibuka',
        message: errorMessage!,
        actionLabel: 'Coba lagi',
        onAction: onRetry,
      );
    }

    if (cameraController == null || !cameraController.value.isInitialized) {
      return const _CameraMessage(
        icon: Icons.camera_alt_outlined,
        title: 'Menyiapkan kamera',
        message: 'Mohon tunggu sebentar.',
        showProgress: true,
      );
    }

    return LayoutBuilder(
      builder: (context, constraints) {
        final previewAspectRatio = cameraController.value.aspectRatio;
        var previewWidth = constraints.maxWidth;
        var previewHeight = previewWidth / previewAspectRatio;

        if (previewHeight < constraints.maxHeight) {
          previewHeight = constraints.maxHeight;
          previewWidth = previewHeight * previewAspectRatio;
        }

        return ClipRect(
          child: Center(
            child: SizedBox(
              width: previewWidth,
              height: previewHeight,
              child: CameraPreview(cameraController),
            ),
          ),
        );
      },
    );
  }
}

class _CameraMessage extends StatelessWidget {
  final IconData icon;
  final String title;
  final String message;
  final bool showProgress;
  final String? actionLabel;
  final VoidCallback? onAction;

  const _CameraMessage({
    required this.icon,
    required this.title,
    required this.message,
    this.showProgress = false,
    this.actionLabel,
    this.onAction,
  });

  @override
  Widget build(BuildContext context) {
    return ColoredBox(
      color: Colors.black,
      child: Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(icon, size: 48, color: AppColors.primaryGreen),
              const SizedBox(height: 16),
              Text(
                title,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  color: Colors.white,
                  fontWeight: FontWeight.w700,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                message,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.white.withValues(alpha: 0.72),
                ),
              ),
              if (showProgress) ...[
                const SizedBox(height: 20),
                const CircularProgressIndicator(color: AppColors.primaryGreen),
              ],
              if (actionLabel != null && onAction != null) ...[
                const SizedBox(height: 20),
                FilledButton(
                  onPressed: onAction,
                  child: Text(actionLabel!),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
