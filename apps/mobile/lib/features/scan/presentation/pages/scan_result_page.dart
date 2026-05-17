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
    final scanState = ref.watch(scanControllerProvider);
    final result = scanState.result;

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
            ScanResultSummaryCard(result: result),
            const SizedBox(height: 16),
            VisualDominanceCard(result: result),
            const SizedBox(height: 16),
            PortionSuggestionCard(result: result),
            const SizedBox(height: 32),
            SizedBox(
              width: double.infinity,
              child: AppButton(
                label: scanState.isSaved
                    ? 'Sudah Tersimpan'
                    : 'Simpan ke Riwayat',
                icon: scanState.isSaved ? Icons.check : Icons.bookmark_add,
                onPressed: result == null || scanState.isSaved
                    ? null
                    : () {
                        ref
                            .read(scanControllerProvider.notifier)
                            .saveCurrentResult();
                        ScaffoldMessenger.of(context)
                          ..hideCurrentSnackBar()
                          ..showSnackBar(
                            const SnackBar(
                              content: Text('Hasil scan tersimpan di riwayat.'),
                            ),
                          );
                        ref.read(scanControllerProvider.notifier).reset();
                        context.go('/history');
                      },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
