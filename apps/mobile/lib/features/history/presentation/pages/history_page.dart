import 'package:flutter/material.dart';
import '../widgets/meal_history_list.dart';

class HistoryPage extends StatelessWidget {
  const HistoryPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Meal History'),
      ),
      body: const MealHistoryList(),
    );
  }
}
