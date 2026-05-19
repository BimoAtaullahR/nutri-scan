import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';
import '../controllers/scan_controller.dart';

class VisualDominanceCard extends StatelessWidget {
  final ScanResult? result;

  const VisualDominanceCard({super.key, this.result});

  @override
  Widget build(BuildContext context) {
    return AppCard(
      color: AppColors.lightBlue.withValues(alpha: 0.1),
      child: Row(
        children: [
          const Icon(Icons.pie_chart, size: 40, color: AppColors.lightBlue),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Dominasi Visual',
                  style: Theme.of(context).textTheme.labelLarge,
                ),
                const SizedBox(height: 4),
                Text(
                  '${result?.dominantPortionLabel ?? 'Karbohidrat'} terlihat dominan, sekitar ${result?.dominantPortionKcal ?? 50} kcal bisa disisihkan.',
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
