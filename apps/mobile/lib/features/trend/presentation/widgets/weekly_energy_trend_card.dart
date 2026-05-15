import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';

class WeeklyEnergyTrendCard extends StatelessWidget {
  const WeeklyEnergyTrendCard({super.key});

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
                'Weekly Energy Trend',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const Icon(Icons.show_chart, color: AppColors.primaryGreen),
            ],
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 150,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                _buildBar(context, 'Sen', 0.6),
                _buildBar(context, 'Sel', 0.8),
                _buildBar(context, 'Rab', 0.5),
                _buildBar(context, 'Kam', 0.9),
                _buildBar(context, 'Jum', 0.7),
                _buildBar(context, 'Sab', 0.4),
                _buildBar(context, 'Min', 0.3),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBar(BuildContext context, String day, double percentage) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        Container(
          width: 24,
          height: 100 * percentage,
          decoration: BoxDecoration(
            color: AppColors.lightBlue,
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          day,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: AppColors.neutralMuted,
              ),
        ),
      ],
    );
  }
}
