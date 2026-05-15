import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class ScannerPreviewFrame extends StatelessWidget {
  const ScannerPreviewFrame({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: 400,
      decoration: BoxDecoration(
        color: AppColors.neutralMuted.withValues(alpha: 0.2),
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: AppColors.primaryGreen, width: 2),
      ),
      child: const Center(
        child: Icon(
          Icons.restaurant,
          size: 64,
          color: AppColors.neutralMuted,
        ),
      ),
    );
  }
}
