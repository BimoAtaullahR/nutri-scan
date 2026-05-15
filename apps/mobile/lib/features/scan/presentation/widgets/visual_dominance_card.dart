import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';

class VisualDominanceCard extends StatelessWidget {
  const VisualDominanceCard({super.key});

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
                  '60% Karbohidrat, 20% Protein, 20% Lemak',
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
