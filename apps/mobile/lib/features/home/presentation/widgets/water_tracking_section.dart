import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class WaterTrackingSection extends StatelessWidget {
  const WaterTrackingSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Track How Much You\'ve Drunk Today',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: AppColors.darkNavy,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 16),
            decoration: BoxDecoration(
              color: AppColors.primaryGreen,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Column(
              children: [
                const Text(
                  'Water',
                  style: TextStyle(
                    color: AppColors.darkNavy,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                Text(
                  'Goal: 2,00 Liters',
                  style: TextStyle(
                    color: AppColors.darkNavy.withValues(alpha: 0.8),
                    fontSize: 12,
                  ),
                ),
                const SizedBox(height: 8),
                const Text(
                  '0,25 Liters',
                  style: TextStyle(
                    color: AppColors.darkNavy,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 16),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: [
                    _buildDrop(active: true),
                    _buildAddDrop(),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                    _buildDrop(active: false),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDrop({required bool active}) {
    return Icon(
      Icons.water_drop,
      color: active
          ? AppColors.lightBlue
          : AppColors.lightBlue.withValues(alpha: 0.5),
      size: 32,
    );
  }

  Widget _buildAddDrop() {
    return Container(
      width: 24,
      height: 24,
      margin: const EdgeInsets.only(top: 4, left: 4, right: 4),
      decoration: BoxDecoration(
        color: AppColors.lightBlue,
        shape: BoxShape.circle,
      ),
      child: const Icon(Icons.add, color: AppColors.darkNavy, size: 16),
    );
  }
}
