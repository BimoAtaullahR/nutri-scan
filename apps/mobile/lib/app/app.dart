import 'package:flutter/material.dart';

import '../features/scan/presentation/pages/scan_page.dart';
import 'theme/app_theme.dart';

class NutriScanApp extends StatelessWidget {
  const NutriScanApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'NutriScan',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: const ScanPage(),
    );
  }
}
