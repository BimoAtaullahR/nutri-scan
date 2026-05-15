import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class AchievementsSection extends StatelessWidget {
  const AchievementsSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Want To Record Your Achievements?',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 8),
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(vertical: 32),
            decoration: BoxDecoration(
              color: AppColors.lightBlue,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                Icon(
                  Icons.insert_drive_file,
                  size: 48,
                  color: AppColors.darkNavy.withValues(alpha: 0.3),
                ),
                const SizedBox(height: 8),
                Text(
                  'Write About Your Achievements Here',
                  style: TextStyle(
                    color: AppColors.darkNavy.withValues(alpha: 0.5),
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
          )
        ],
      ),
    );
  }
}
