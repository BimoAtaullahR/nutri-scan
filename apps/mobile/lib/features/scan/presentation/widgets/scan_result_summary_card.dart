import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';
import '../controllers/scan_controller.dart';

class ScanResultSummaryCard extends StatelessWidget {
  final ScanResult? result;

  const ScanResultSummaryCard({super.key, this.result});

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                result?.foodName ?? 'Nasi Goreng Ayam',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const Icon(Icons.check_circle, color: AppColors.primaryGreen),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            'Estimasi Energi: ${result?.estimatedEnergyKcal ?? 450} kcal',
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: AppColors.energyOrange,
            ),
          ),
        ],
      ),
    );
  }
}
