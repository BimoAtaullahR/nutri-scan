import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class ProgressSection extends StatelessWidget {
  const ProgressSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Look At Your Progress Today!',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 8),
          const _CalendarStrip(),
          const SizedBox(height: 8),
          Align(
            alignment: Alignment.centerRight,
            child: Text(
              'See Details',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: AppColors.darkNavy,
                    fontWeight: FontWeight.bold,
                  ),
            ),
          ),
          const SizedBox(height: 8),
          const _ProgressCard(),
          const SizedBox(height: 12),
          const _MealCard(title: 'Breakfast', kcal: '0/552', icon: Icons.breakfast_dining),
          const _MealCard(title: 'Lunch', kcal: '0/552', icon: Icons.lunch_dining),
          const _MealCard(title: 'Dinner', kcal: '0/552', icon: Icons.dinner_dining),
          const _MealCard(title: 'Snack', kcal: '0/552', icon: Icons.apple),
        ],
      ),
    );
  }
}

class _CalendarStrip extends StatelessWidget {
  const _CalendarStrip();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.primaryGreen,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'May | 2026',
            style: Theme.of(context).textTheme.titleSmall?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _buildDay(context, 'S', '01'),
              _buildDay(context, 'M', '02'),
              _buildDay(context, 'T', '03'),
              _buildDay(context, 'W', '04'),
              _buildDay(context, 'T', '05', isSelected: true),
              _buildDay(context, 'F', '06'),
              _buildDay(context, 'S', '07'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDay(BuildContext context, String day, String date, {bool isSelected = false}) {
    return Column(
      children: [
        Text(
          day,
          style: TextStyle(
            color: isSelected ? AppColors.lightBlue : AppColors.lightBlue,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 8),
          decoration: BoxDecoration(
            color: isSelected ? AppColors.darkNavy : Colors.transparent,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Text(
            date,
            style: TextStyle(
              color: isSelected ? Colors.white : AppColors.lightBlue,
              fontWeight: FontWeight.bold,
              fontSize: 18,
            ),
          ),
        ),
      ],
    );
  }
}

class _ProgressCard extends StatelessWidget {
  const _ProgressCard();

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.lightBlue,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          const SizedBox(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _buildStat('0', 'Eaten'),
              _buildArcProgress(),
              _buildStat('0', 'Burned'),
            ],
          ),
          const SizedBox(height: 24),
          Container(
            padding: const EdgeInsets.symmetric(vertical: 12),
            decoration: const BoxDecoration(
              color: Color(0xFFC7D59F), // Darker lime green
              borderRadius: BorderRadius.only(
                bottomLeft: Radius.circular(16),
                bottomRight: Radius.circular(16),
              ),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                _buildMacro('Carbs', '0/244 g'),
                _buildMacro('Protein', '0/90 g'),
                _buildMacro('Fat', '0/59 g'),
              ],
            ),
          )
        ],
      ),
    );
  }

  Widget _buildStat(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(color: AppColors.darkNavy, fontSize: 20, fontWeight: FontWeight.bold),
        ),
        Text(
          label,
          style: const TextStyle(color: AppColors.darkNavy, fontSize: 12),
        ),
      ],
    );
  }

  Widget _buildArcProgress() {
    return Stack(
      alignment: Alignment.center,
      children: [
        SizedBox(
          width: 120,
          height: 120,
          child: CircularProgressIndicator(
            value: 0.8, // Example
            strokeWidth: 16,
            backgroundColor: AppColors.darkNavy.withValues(alpha: 0.1),
            valueColor: AlwaysStoppedAnimation<Color>(AppColors.darkNavy.withValues(alpha: 0.3)),
          ),
        ),
        const Column(
          children: [
            Text(
              '1.836',
              style: TextStyle(color: AppColors.darkNavy, fontSize: 24, fontWeight: FontWeight.bold),
            ),
            Text(
              'Remaining',
              style: TextStyle(color: AppColors.darkNavy, fontSize: 14),
            ),
          ],
        )
      ],
    );
  }

  Widget _buildMacro(String label, String value) {
    return Column(
      children: [
        Text(
          label,
          style: const TextStyle(color: AppColors.darkNavy, fontSize: 12, fontWeight: FontWeight.w500),
        ),
        Text(
          value,
          style: const TextStyle(color: AppColors.darkNavy, fontSize: 14, fontWeight: FontWeight.w500),
        ),
      ],
    );
  }
}

class _MealCard extends StatelessWidget {
  final String title;
  final String kcal;
  final IconData icon;

  const _MealCard({required this.title, required this.kcal, required this.icon});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 4),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: AppColors.lightBlue.withValues(alpha: 0.2),
        border: Border.all(color: AppColors.lightBlue.withValues(alpha: 0.3)),
        borderRadius: title == 'Breakfast' 
          ? const BorderRadius.only(topLeft: Radius.circular(12), topRight: Radius.circular(12))
          : title == 'Snack'
            ? const BorderRadius.only(bottomLeft: Radius.circular(12), bottomRight: Radius.circular(12))
            : BorderRadius.zero,
      ),
      child: Row(
        children: [
          CircleAvatar(
            backgroundColor: AppColors.darkNavy,
            radius: 16,
            child: Icon(icon, color: Colors.white, size: 16),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(title, style: const TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.bold)),
                    const SizedBox(width: 8),
                    const Text('>', style: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.bold)),
                  ],
                ),
                Text('$kcal kcal', style: TextStyle(color: AppColors.darkNavy.withValues(alpha: 0.7), fontSize: 12)),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(4),
            decoration: const BoxDecoration(
              color: AppColors.darkNavy,
              shape: BoxShape.circle,
            ),
            child: const Icon(Icons.add, color: Colors.white, size: 16),
          )
        ],
      ),
    );
  }
}
