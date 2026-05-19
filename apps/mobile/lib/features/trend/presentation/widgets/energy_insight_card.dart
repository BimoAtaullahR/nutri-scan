import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';

class EnergyInsightCard extends StatelessWidget {
  const EnergyInsightCard({super.key});

  @override
  Widget build(BuildContext context) {
    return AppCard(
      color: AppColors.mintSurface,
      child: Row(
        children: [
          const Icon(Icons.auto_graph, color: AppColors.darkNavy, size: 32),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Trend Stabil',
                  style: Theme.of(context).textTheme.labelLarge,
                ),
                const SizedBox(height: 4),
                Text(
                  'Trend energi mingguan Anda terpantau stabil. Lanjutkan kebiasaan baik ini!',
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
