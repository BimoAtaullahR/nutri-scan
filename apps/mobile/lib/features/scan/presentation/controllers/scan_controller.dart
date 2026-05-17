import 'package:flutter_riverpod/flutter_riverpod.dart';

enum ScanStatus { idle, analyzing, success, error }

class ScanResult {
  final String foodName;
  final int estimatedEnergyKcal;
  final String auraHeadline;
  final String auraSuggestion;
  final String dominantPortionLabel;
  final int dominantPortionKcal;
  final bool suggestionFollowed;
  final DateTime capturedAt;

  const ScanResult({
    required this.foodName,
    required this.estimatedEnergyKcal,
    required this.auraHeadline,
    required this.auraSuggestion,
    required this.dominantPortionLabel,
    required this.dominantPortionKcal,
    this.suggestionFollowed = false,
    required this.capturedAt,
  });

  ScanResult copyWith({
    String? foodName,
    int? estimatedEnergyKcal,
    String? auraHeadline,
    String? auraSuggestion,
    String? dominantPortionLabel,
    int? dominantPortionKcal,
    bool? suggestionFollowed,
    DateTime? capturedAt,
  }) {
    return ScanResult(
      foodName: foodName ?? this.foodName,
      estimatedEnergyKcal: estimatedEnergyKcal ?? this.estimatedEnergyKcal,
      auraHeadline: auraHeadline ?? this.auraHeadline,
      auraSuggestion: auraSuggestion ?? this.auraSuggestion,
      dominantPortionLabel: dominantPortionLabel ?? this.dominantPortionLabel,
      dominantPortionKcal: dominantPortionKcal ?? this.dominantPortionKcal,
      suggestionFollowed: suggestionFollowed ?? this.suggestionFollowed,
      capturedAt: capturedAt ?? this.capturedAt,
    );
  }
}

class ScanState {
  final ScanStatus status;
  final String? errorMessage;
  final String? capturedImagePath;
  final ScanResult? result;
  final bool isAuraPlateRevealed;
  final bool isSaved;

  ScanState({
    this.status = ScanStatus.idle,
    this.errorMessage,
    this.capturedImagePath,
    this.result,
    this.isAuraPlateRevealed = false,
    this.isSaved = false,
  });

  ScanState copyWith({
    ScanStatus? status,
    String? errorMessage,
    String? capturedImagePath,
    ScanResult? result,
    bool? isAuraPlateRevealed,
    bool? isSaved,
  }) {
    return ScanState(
      status: status ?? this.status,
      errorMessage: errorMessage ?? this.errorMessage,
      capturedImagePath: capturedImagePath ?? this.capturedImagePath,
      result: result ?? this.result,
      isAuraPlateRevealed: isAuraPlateRevealed ?? this.isAuraPlateRevealed,
      isSaved: isSaved ?? this.isSaved,
    );
  }
}

class ScanController extends Notifier<ScanState> {
  @override
  ScanState build() {
    return ScanState();
  }

  void reset() {
    state = ScanState();
  }

  void fail(String message) {
    state = ScanState(status: ScanStatus.error, errorMessage: message);
  }

  void toggleAuraPlate() {
    if (state.status != ScanStatus.success) return;

    state = state.copyWith(isAuraPlateRevealed: !state.isAuraPlateRevealed);
  }

  void saveCurrentResult() {
    final result = state.result;
    if (result == null || state.isSaved) return;

    ref.read(savedScanResultsProvider.notifier).save(result);
    state = state.copyWith(isSaved: true);
  }

  void setSuggestionFollowed(bool value) {
    final result = state.result;
    if (result == null) return;

    state = state.copyWith(result: result.copyWith(suggestionFollowed: value));
  }

  Future<void> analyzeImage({required String capturedImagePath}) async {
    state = ScanState(
      status: ScanStatus.analyzing,
      capturedImagePath: capturedImagePath,
    );
    await Future.delayed(const Duration(seconds: 2));
    state = state.copyWith(
      status: ScanStatus.success,
      result: ScanResult(
        foodName: 'Salmon Rice Bowl',
        estimatedEnergyKcal: 580,
        auraHeadline: 'Piringmu terlihat seimbang!',
        auraSuggestion:
            'Sisihkan sedikit nasi untuk menghemat sekitar 50 kcal.',
        dominantPortionLabel: 'Karbohidrat',
        dominantPortionKcal: 50,
        suggestionFollowed: false,
        capturedAt: DateTime.now(),
      ),
    );
  }
}

final scanControllerProvider = NotifierProvider<ScanController, ScanState>(() {
  return ScanController();
});

class SavedScanResultsController extends Notifier<List<ScanResult>> {
  @override
  List<ScanResult> build() {
    return const [];
  }

  void save(ScanResult result) {
    state = [result.copyWith(capturedAt: DateTime.now()), ...state];
  }
}

final savedScanResultsProvider =
    NotifierProvider<SavedScanResultsController, List<ScanResult>>(() {
      return SavedScanResultsController();
    });
