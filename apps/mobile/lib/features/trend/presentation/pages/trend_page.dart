import 'package:flutter/material.dart';
import '../widgets/weekly_energy_trend_card.dart';
import '../widgets/energy_insight_card.dart';

class TrendPage extends StatelessWidget {
  const TrendPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Energy Trend')),
      body: const SingleChildScrollView(
        padding: EdgeInsets.all(24.0),
        child: Column(
          children: [
            WeeklyEnergyTrendCard(),
            SizedBox(height: 16),
            EnergyInsightCard(),
            SizedBox(height: 32),
          ],
        ),
      ),
    );
  }
}
