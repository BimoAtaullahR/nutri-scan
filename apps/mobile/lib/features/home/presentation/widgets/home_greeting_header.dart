import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class HomeGreetingHeader extends StatelessWidget {
  const HomeGreetingHeader({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Hello,',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.w900,
                ),
              ),
              Text(
                'Daniel Chald',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  color: AppColors.darkNavy,
                  fontWeight: FontWeight.w400,
                ),
              ),
            ],
          ),
          const CircleAvatar(
            radius: 20,
            backgroundColor: AppColors.primaryGreen,
            child: Icon(Icons.person, color: AppColors.darkNavy, size: 24),
          ),
        ],
      ),
    );
  }
}
