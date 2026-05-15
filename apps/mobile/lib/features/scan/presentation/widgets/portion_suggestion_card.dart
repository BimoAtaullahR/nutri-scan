import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';

class PortionSuggestionCard extends StatelessWidget {
  const PortionSuggestionCard({super.key});

  @override
  Widget build(BuildContext context) {
    return AppCard(
      color: AppColors.primaryGreen.withValues(alpha: 0.1),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.lightbulb_outline, color: AppColors.primaryGreen),
              const SizedBox(width: 8),
              Text(
                'Saran Porsi',
                style: Theme.of(context).textTheme.titleMedium,
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            'Kurangi sedikit nasi sekitar 2 sendok makan untuk menyeimbangkan asupan kalori Anda hari ini.',
            style: Theme.of(context).textTheme.bodyMedium,
          ),
        ],
      ),
    );
  }
}
