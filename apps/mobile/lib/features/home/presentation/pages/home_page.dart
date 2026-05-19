import 'package:flutter/material.dart';
import '../widgets/home_greeting_header.dart';
import '../widgets/scan_food_cta_card.dart';
import '../widgets/progress_section.dart';
import '../widgets/water_tracking_section.dart';
import '../widgets/achievements_section.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: 16),
              HomeGreetingHeader(),
              SizedBox(height: 16),
              ScanFoodCtaCard(),
              SizedBox(height: 24),
              ProgressSection(),
              SizedBox(height: 24),
              WaterTrackingSection(),
              SizedBox(height: 24),
              AchievementsSection(),
              SizedBox(height: 48),
            ],
          ),
        ),
      ),
    );
  }
}
