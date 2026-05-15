import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../widgets/scan_result_summary_card.dart';
import '../widgets/visual_dominance_card.dart';
import '../widgets/portion_suggestion_card.dart';
import '../../../../shared/widgets/app_button.dart';
import '../controllers/scan_controller.dart';

class ScanResultPage extends ConsumerWidget {
  const ScanResultPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Hasil Scan'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () {
            ref.read(scanControllerProvider.notifier).reset();
            context.pop();
          },
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const ScanResultSummaryCard(),
            const SizedBox(height: 16),
            const VisualDominanceCard(),
            const SizedBox(height: 16),
            const PortionSuggestionCard(),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              child: AppButton(
                label: 'Simpan ke Riwayat',
                onPressed: () {
                  ref.read(scanControllerProvider.notifier).reset();
                  context.go('/');
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
