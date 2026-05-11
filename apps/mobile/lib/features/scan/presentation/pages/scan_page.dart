import 'package:flutter/material.dart';

class ScanPage extends StatelessWidget {
  const ScanPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('NutriScan'),
      ),
      body: const Center(
        child: Text('AI Vision Scanner Page'),
      ),
    );
  }
}
