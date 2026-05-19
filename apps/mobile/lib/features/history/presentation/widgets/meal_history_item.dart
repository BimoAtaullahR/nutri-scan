import 'package:flutter/material.dart';

import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_card.dart';

class MealHistoryItem extends StatelessWidget {
  final String title;
  final String calories;
  final String time;
  final IconData icon;
  final String? note;

  const MealHistoryItem({
    super.key,
    required this.title,
    required this.calories,
    required this.time,
    required this.icon,
    this.note,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: AppCard(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: AppColors.lightBlue.withValues(alpha: 0.2),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, color: AppColors.lightBlue),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: Theme.of(context).textTheme.labelLarge),
                  const SizedBox(height: 4),
                  Text(
                    '$calories kcal - $time',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: AppColors.neutralMuted,
                    ),
                  ),
                  if (note != null) ...[
                    const SizedBox(height: 4),
                    Text(
                      note!,
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppColors.neutralBody,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ],
              ),
            ),
            const Icon(Icons.chevron_right, color: AppColors.neutralMuted),
          ],
        ),
      ),
    );
  }
}
