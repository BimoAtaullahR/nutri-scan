import 'package:flutter/material.dart';
import 'meal_history_item.dart';

class MealHistoryList extends StatelessWidget {
  const MealHistoryList({super.key});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(24),
      children: const [
        Text(
          'Hari Ini',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        SizedBox(height: 12),
        MealHistoryItem(
          title: 'Nasi Goreng Ayam',
          calories: '450',
          time: '12:30 PM',
          icon: Icons.fastfood,
        ),
        MealHistoryItem(
          title: 'Roti Gandum & Telur',
          calories: '300',
          time: '08:00 AM',
          icon: Icons.breakfast_dining,
        ),
        SizedBox(height: 24),
        Text(
          'Kemarin',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
        ),
        SizedBox(height: 12),
        MealHistoryItem(
          title: 'Ayam Bakar Dada',
          calories: '350',
          time: '07:00 PM',
          icon: Icons.dinner_dining,
        ),
        MealHistoryItem(
          title: 'Salad Buah',
          calories: '150',
          time: '04:00 PM',
          icon: Icons.local_dining,
        ),
      ],
    );
  }
}
