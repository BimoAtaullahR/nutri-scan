import 'package:flutter/material.dart';

class AppColors {
  // Theme from Figma / Images
  static const Color primaryGreen = Color(0xFFE7F1A8); // Lime Green
  static const Color darkNavy = Color(0xFF364C84); // Typography & Authority
  static const Color lightBlue = Color(0xFF95B1EE); // The Magical Vibe
  static const Color backgroundCream = Color(0xFFFFFDF5); // Warm Cream

  // Supporting Colors
  static const Color cardWhite = Colors.white;
  static const Color neutralBody = Color(0xFF4B5563);
  static const Color neutralMuted = Color(0xFF9CA3AF);
  static const Color borderSoft = Color(0xFFE5E7EB);
  static const Color energyOrange = Color(0xFFF97316);
  static const Color softOrange = Color(0xFFFFEDD5);
  static const Color mintSurface = Color(0xFFECFDF5);
}

class AppTheme {
  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      scaffoldBackgroundColor: AppColors.backgroundCream,
      colorScheme: const ColorScheme.light(
        primary: AppColors.primaryGreen,
        secondary: AppColors.lightBlue,
        surface: AppColors.backgroundCream,
        onSurface: AppColors.darkNavy,
        onPrimary: AppColors.darkNavy,
      ),
      textTheme: const TextTheme(
        displayLarge: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.bold),
        displayMedium: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.bold),
        titleLarge: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.w700),
        titleMedium: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.w600),
        bodyLarge: TextStyle(color: AppColors.neutralBody),
        bodyMedium: TextStyle(color: AppColors.neutralBody),
        labelLarge: TextStyle(color: AppColors.darkNavy, fontWeight: FontWeight.w600),
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.backgroundCream,
        elevation: 0,
        centerTitle: false,
        iconTheme: IconThemeData(color: AppColors.darkNavy),
        titleTextStyle: TextStyle(
          color: AppColors.darkNavy,
          fontSize: 20,
          fontWeight: FontWeight.bold,
        ),
      ),
      cardTheme: CardThemeData(
        color: AppColors.cardWhite,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(24),
          side: const BorderSide(color: AppColors.borderSoft, width: 1),
        ),
        margin: EdgeInsets.zero,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.primaryGreen,
          foregroundColor: AppColors.darkNavy,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
          textStyle: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(24),
          ),
        ),
      ),
    );
  }
}
