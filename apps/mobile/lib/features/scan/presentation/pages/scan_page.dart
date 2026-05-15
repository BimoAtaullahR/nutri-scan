import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../controllers/scan_controller.dart';
import '../widgets/scanner_preview_frame.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../app/theme/app_theme.dart';

class ScanPage extends ConsumerWidget {
  const ScanPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scanState = ref.watch(scanControllerProvider);

    // Listen to state changes to navigate when success
    ref.listen<ScanState>(scanControllerProvider, (previous, next) {
      if (next.status == ScanStatus.success) {
        context.pushReplacement('/result');
      }
    });

    return Scaffold(
      appBar: AppBar(
        title: const Text('Scan Food'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => context.pop(),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          children: [
            const Text(
              'Arahkan kamera ke makanan Anda',
              style: TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 24),
            const ScannerPreviewFrame(),
            const Spacer(),
            if (scanState.status == ScanStatus.analyzing)
              const Column(
                children: [
                  CircularProgressIndicator(color: AppColors.primaryGreen),
                  SizedBox(height: 16),
                  Text('Menganalisis nutrisi...'),
                ],
              )
            else
              SizedBox(
                width: double.infinity,
                child: AppButton(
                  label: 'Capture & Analyze',
                  icon: Icons.camera,
                  onPressed: () {
                    ref.read(scanControllerProvider.notifier).analyzeImage();
                  },
                ),
              ),
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }
}
