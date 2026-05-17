import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../app/theme/app_theme.dart';
import '../../../../shared/widgets/app_button.dart';
import '../../../../shared/widgets/app_card.dart';
import '../controllers/scan_controller.dart';
import '../widgets/portion_suggestion_card.dart';
import '../widgets/scan_result_summary_card.dart';
import '../widgets/visual_dominance_card.dart';

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
            context.go('/');
          },
        ),
      ),
      body: result == null
          ? const _EmptyResult()
          : SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(24, 12, 24, 32),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  ScanResultSummaryCard(result: result),
                  const SizedBox(height: 16),
                  _NutritionDetailCard(result: result),
                  const SizedBox(height: 16),
                  VisualDominanceCard(result: result),
                  const SizedBox(height: 16),
                  PortionSuggestionCard(result: result),
                  const SizedBox(height: 16),
                  _SuggestionConfirmationCard(
                    result: result,
                    onChanged: (value) => ref
                        .read(scanControllerProvider.notifier)
                        .setSuggestionFollowed(value ?? false),
                  ),
                  const SizedBox(height: 28),
                  SizedBox(
                    width: double.infinity,
                    child: AppButton(
                      label: scanState.isSaved ? 'Saved' : 'Save and Finish',
                      icon: scanState.isSaved ? Icons.check : Icons.task_alt,
                      onPressed: scanState.isSaved
                          ? null
                          : () {
                              ref
                                  .read(scanControllerProvider.notifier)
                                  .saveCurrentResult();
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

class _NutritionDetailCard extends StatelessWidget {
  final ScanResult result;

  const _NutritionDetailCard({required this.result});

  @override
  Widget build(BuildContext context) {
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.monitor_heart, color: AppColors.energyOrange),
              const SizedBox(width: 8),
              Text(
                'Detail Nutrisi',
                style: Theme.of(context).textTheme.titleMedium,
              ),
            ],
          ),
          const SizedBox(height: 14),
          _NutritionRow(
            label: 'Total estimasi',
            value: '~${result.estimatedEnergyKcal} kcal',
          ),
          _NutritionRow(
            label: 'Porsi dominan',
            value: result.dominantPortionLabel,
          ),
          _NutritionRow(
            label: 'Saran disisihkan',
            value: '~${result.dominantPortionKcal} kcal',
          ),
          const _NutritionRow(
            label: 'Basis estimasi',
            value: 'Visual piring dan komposisi makanan',
          ),
        ],
      ),
    );
  }
}

class _NutritionRow extends StatelessWidget {
  final String label;
  final String value;

  const _NutritionRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Text(
              label,
              style: Theme.of(
                context,
              ).textTheme.bodyMedium?.copyWith(color: AppColors.neutralMuted),
            ),
          ),
          const SizedBox(width: 16),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.right,
              style: Theme.of(context).textTheme.labelLarge,
            ),
          ),
        ],
      ),
    );
  }
}

class _SuggestionConfirmationCard extends StatelessWidget {
  final ScanResult result;
  final ValueChanged<bool?> onChanged;

  const _SuggestionConfirmationCard({
    required this.result,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return AppCard(
      color: AppColors.mintSurface,
      padding: EdgeInsets.zero,
      child: CheckboxListTile(
        value: result.suggestionFollowed,
        onChanged: onChanged,
        activeColor: AppColors.darkNavy,
        checkboxShape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(4),
        ),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        title: Text(
          'Saya sudah mengikuti saran AI',
          style: Theme.of(context).textTheme.labelLarge,
        ),
        subtitle: Text(
          'Konfirmasi ini akan ikut tersimpan bersama hasil scan.',
          style: Theme.of(
            context,
          ).textTheme.bodySmall?.copyWith(color: AppColors.neutralBody),
        ),
        controlAffinity: ListTileControlAffinity.leading,
      ),
    );
  }
}

class _EmptyResult extends ConsumerWidget {
  const _EmptyResult();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.document_scanner_outlined,
              size: 48,
              color: AppColors.neutralMuted,
            ),
            const SizedBox(height: 16),
            Text(
              'Belum ada hasil scan',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              'Ambil foto makanan dulu untuk melihat estimasi nutrisi.',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            const SizedBox(height: 24),
            AppButton(
              label: 'Kembali ke Home',
              onPressed: () {
                ref.read(scanControllerProvider.notifier).reset();
                context.go('/');
              },
            ),
          ],
        ),
      ),
    );
  }
}
