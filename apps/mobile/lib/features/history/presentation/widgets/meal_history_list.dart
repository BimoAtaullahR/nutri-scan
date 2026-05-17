import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../scan/presentation/controllers/scan_controller.dart';
import 'meal_history_item.dart';

class MealHistoryList extends ConsumerWidget {
  const MealHistoryList({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final savedResults = ref.watch(savedScanResultsProvider);

    return ListView(
      padding: const EdgeInsets.all(24),
      children: [
        if (savedResults.isNotEmpty) ...[
          const Text(
            'Tersimpan',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          for (final result in savedResults)
            MealHistoryItem(
              title: result.foodName,
              calories: result.estimatedEnergyKcal.toString(),
              time: _formatTime(result.capturedAt),
              icon: Icons.restaurant_menu,
            ),
          const SizedBox(height: 24),
        ],
        const Text(
          'Hari Ini',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        const MealHistoryItem(
          title: 'Nasi Goreng Ayam',
          calories: '450',
          time: '12:30 PM',
          icon: Icons.fastfood,
        ),
        const MealHistoryItem(
          title: 'Roti Gandum & Telur',
          calories: '300',
          time: '08:00 AM',
          icon: Icons.breakfast_dining,
        ),
        const SizedBox(height: 24),
        const Text(
          'Kemarin',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        const SizedBox(height: 12),
        const MealHistoryItem(
          title: 'Ayam Bakar Dada',
          calories: '350',
          time: '07:00 PM',
          icon: Icons.dinner_dining,
        ),
        const MealHistoryItem(
          title: 'Salad Buah',
          calories: '150',
          time: '04:00 PM',
          icon: Icons.local_dining,
        ),
      ],
    );
  }

  String _formatTime(DateTime value) {
    final hour = value.hour.toString().padLeft(2, '0');
    final minute = value.minute.toString().padLeft(2, '0');
    return '$hour:$minute';
  }
}
