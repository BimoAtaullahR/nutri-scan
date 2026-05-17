import 'package:flutter_riverpod/flutter_riverpod.dart';

enum ScanStatus { idle, analyzing, success, error }

class ScanState {
  final ScanStatus status;
  final String? errorMessage;
  final String? capturedImagePath;

  ScanState({
    this.status = ScanStatus.idle,
    this.errorMessage,
    this.capturedImagePath,
  });

  ScanState copyWith({
    ScanStatus? status,
    String? errorMessage,
    String? capturedImagePath,
  }) {
    return ScanState(
      status: status ?? this.status,
      errorMessage: errorMessage,
      capturedImagePath: capturedImagePath ?? this.capturedImagePath,
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

  Future<void> analyzeImage({required String capturedImagePath}) async {
    state = ScanState(
      status: ScanStatus.analyzing,
      capturedImagePath: capturedImagePath,
    );
    await Future.delayed(const Duration(seconds: 2));
    state = state.copyWith(status: ScanStatus.success);
  }
}

final scanControllerProvider = NotifierProvider<ScanController, ScanState>(() {
  return ScanController();
});
