import 'dart:math' as math;
import 'package:flutter/material.dart';
import '../../../../app/theme/app_theme.dart';

class ProgressData {
  final int eaten;
  final int burned;
  final int remaining;
  final int totalGoal;
  final double carbsConsumed;
  final double carbsGoal;
  final double proteinConsumed;
  final double proteinGoal;
  final double fatConsumed;
  final double fatGoal;

  const ProgressData({
    required this.eaten,
    required this.burned,
    required this.remaining,
    required this.totalGoal,
    required this.carbsConsumed,
    required this.carbsGoal,
    required this.proteinConsumed,
    required this.proteinGoal,
    required this.fatConsumed,
    required this.fatGoal,
  });

  static const ProgressData dummy = ProgressData(
    eaten: 0,
    burned: 0,
    remaining: 1836,
    totalGoal: 1836,
    carbsConsumed: 0,
    carbsGoal: 244,
    proteinConsumed: 0,
    proteinGoal: 90,
    fatConsumed: 0,
    fatGoal: 59,
  );
}

class MealData {
  final String title;
  final int consumed;
  final int goal;
  final IconData icon;
  final Color color;

  const MealData({
    required this.title,
    required this.consumed,
    required this.goal,
    required this.icon,
    required this.color,
  });

  static const List<MealData> dummyList = [
    MealData(
      title: 'Breakfast',
      consumed: 0,
      goal: 552,
      icon: Icons.egg_alt_rounded,
      color: Color(0xFFF1F1E5),
    ),
    MealData(
      title: 'Lunch',
      consumed: 0,
      goal: 552,
      icon: Icons.lunch_dining_rounded,
      color: Color(0xFF9FBAF1),
    ),
    MealData(
      title: 'Dinner',
      consumed: 0,
      goal: 552,
      icon: Icons.room_service_rounded,
      color: Color(0xFF89A6E0),
    ),
    MealData(
      title: 'Snack',
      consumed: 0,
      goal: 180,
      icon: Icons.apple_rounded,
      color: Color(0xFF7396DC),
    ),
  ];
}

class ProgressSection extends StatelessWidget {
  final ProgressData progressData;
  final List<MealData> meals;

  const ProgressSection({
    super.key,
    this.progressData = ProgressData.dummy,
    List<MealData>? meals,
  }) : meals = meals ?? MealData.dummyList;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24),
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: const Color(0xFFDEDED2),
          borderRadius: BorderRadius.circular(24),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Look At Your Progress Today!',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                color: AppColors.darkNavy,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            const _CalendarStrip(),
            Align(
              alignment: Alignment.centerRight,
              child: Padding(
                padding: const EdgeInsets.only(top: 8, bottom: 8),
                child: GestureDetector(
                  onTap: () {},
                  child: Text(
                    'See Details',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: AppColors.darkNavy,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
            ),
            _ProgressCard(data: progressData),
            const SizedBox(height: 16),

            _MealCardGroup(meals: meals),
          ],
        ),
      ),
    );
  }
}

class _CalendarStrip extends StatefulWidget {
  const _CalendarStrip();

  @override
  State<_CalendarStrip> createState() => _CalendarStripState();
}

class _CalendarStripState extends State<_CalendarStrip> {
  late DateTime _selectedDate;
  late DateTime _weekStart;

  @override
  void initState() {
    super.initState();
    final today = DateTime.now();
    _selectedDate = today;
    _weekStart = today.subtract(Duration(days: today.weekday % 7));
  }

  void _onDayTapped(DateTime date) {
    setState(() => _selectedDate = date);
  }

  void _openCalendarModal() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => _FullCalendarModal(
        selectedDate: _selectedDate,
        onDateSelected: (date) {
          setState(() {
            _selectedDate = date;
            _weekStart = date.subtract(Duration(days: date.weekday % 7));
          });
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final monthLabel =
        '${_monthName(_selectedDate.month)} | ${_selectedDate.year}';

    final days = List.generate(7, (i) => _weekStart.add(Duration(days: i)));

    return GestureDetector(
      onTap: _openCalendarModal,
      child: Container(
        padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
        decoration: BoxDecoration(
          color: AppColors.primaryGreen,
          borderRadius: BorderRadius.circular(16),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  monthLabel,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    color: AppColors.darkNavy,
                    fontWeight: FontWeight.bold,
                    fontSize: 15,
                  ),
                ),
                const Spacer(),
                const Icon(
                  Icons.calendar_month_rounded,
                  size: 18,
                  color: AppColors.darkNavy,
                ),
              ],
            ),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: days.map((date) {
                final isSelected = _isSameDay(date, _selectedDate);
                return _DayCell(
                  dayLetter: _dayLetter(date.weekday),
                  date: date.day,
                  isSelected: isSelected,
                  onTap: () => _onDayTapped(date),
                );
              }).toList(),
            ),
          ],
        ),
      ),
    );
  }

  static bool _isSameDay(DateTime a, DateTime b) =>
      a.year == b.year && a.month == b.month && a.day == b.day;

  static String _dayLetter(int weekday) {
    const letters = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];
    if (weekday == 7) return 'S';
    if (weekday == 1) return 'M';
    if (weekday == 2) return 'T';
    if (weekday == 3) return 'W';
    if (weekday == 4) return 'T';
    if (weekday == 5) return 'F';
    return 'S';
  }

  static String _monthName(int month) {
    const names = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec',
    ];
    return names[month - 1];
  }
}

class _DayCell extends StatelessWidget {
  final String dayLetter;
  final int date;
  final bool isSelected;
  final VoidCallback onTap;

  const _DayCell({
    required this.dayLetter,
    required this.date,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 180),
        curve: Curves.easeInOut,
        width: 38,
        padding: const EdgeInsets.symmetric(vertical: 6),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.darkNavy : Colors.transparent,
          borderRadius: BorderRadius.circular(10),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              dayLetter,
              style: TextStyle(
                color: isSelected ? Colors.white : AppColors.lightBlue,
                fontWeight: FontWeight.w600,
                fontSize: 16,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              date.toString().padLeft(2, '0'),
              style: TextStyle(
                color: isSelected ? Colors.white : AppColors.lightBlue,
                fontWeight: FontWeight.bold,
                fontSize: 20,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _FullCalendarModal extends StatefulWidget {
  final DateTime selectedDate;
  final ValueChanged<DateTime> onDateSelected;

  const _FullCalendarModal({
    required this.selectedDate,
    required this.onDateSelected,
  });

  @override
  State<_FullCalendarModal> createState() => _FullCalendarModalState();
}

class _FullCalendarModalState extends State<_FullCalendarModal> {
  late DateTime _viewMonth;
  late DateTime _selected;

  @override
  void initState() {
    super.initState();
    _selected = widget.selectedDate;
    _viewMonth = DateTime(_selected.year, _selected.month);
  }

  void _prevMonth() {
    setState(
      () => _viewMonth = DateTime(_viewMonth.year, _viewMonth.month - 1),
    );
  }

  void _nextMonth() {
    setState(
      () => _viewMonth = DateTime(_viewMonth.year, _viewMonth.month + 1),
    );
  }

  @override
  Widget build(BuildContext context) {
    final daysInMonth = DateUtils.getDaysInMonth(
      _viewMonth.year,
      _viewMonth.month,
    );
    final firstWeekday = DateTime(_viewMonth.year, _viewMonth.month, 1).weekday;
    final offset = firstWeekday % 7;

    return DraggableScrollableSheet(
      initialChildSize: 0.60,
      minChildSize: 0.45,
      maxChildSize: 0.85,
      expand: false,
      builder: (_, controller) => Container(
        decoration: const BoxDecoration(
          color: AppColors.backgroundCream,
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          children: [
            Container(
              margin: const EdgeInsets.only(top: 10),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.darkNavy.withValues(alpha: 0.2),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 16, 20, 8),
              child: Row(
                children: [
                  IconButton(
                    onPressed: _prevMonth,
                    icon: const Icon(
                      Icons.chevron_left_rounded,
                      color: AppColors.darkNavy,
                    ),
                  ),
                  Expanded(
                    child: Text(
                      '${_monthName(_viewMonth.month)} ${_viewMonth.year}',
                      textAlign: TextAlign.center,
                      style: const TextStyle(
                        color: AppColors.darkNavy,
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                  ),
                  IconButton(
                    onPressed: _nextMonth,
                    icon: const Icon(
                      Icons.chevron_right_rounded,
                      color: AppColors.darkNavy,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: const ['S', 'M', 'T', 'W', 'T', 'F', 'S']
                    .map(
                      (d) => SizedBox(
                        width: 36,
                        child: Text(
                          d,
                          textAlign: TextAlign.center,
                          style: const TextStyle(
                            color: AppColors.darkNavy,
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                          ),
                        ),
                      ),
                    )
                    .toList(),
              ),
            ),
            const SizedBox(height: 8),
            Expanded(
              child: SingleChildScrollView(
                controller: controller,
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: GridView.builder(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 7,
                          mainAxisSpacing: 4,
                          crossAxisSpacing: 0,
                          childAspectRatio: 1,
                        ),
                    itemCount: offset + daysInMonth,
                    itemBuilder: (_, index) {
                      if (index < offset) return const SizedBox.shrink();
                      final day = index - offset + 1;
                      final date = DateTime(
                        _viewMonth.year,
                        _viewMonth.month,
                        day,
                      );
                      final isSelected =
                          date.year == _selected.year &&
                          date.month == _selected.month &&
                          date.day == _selected.day;
                      final isToday = DateUtils.isSameDay(date, DateTime.now());

                      return GestureDetector(
                        onTap: () {
                          setState(() => _selected = date);
                          widget.onDateSelected(date);
                          Navigator.pop(context);
                        },
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 150),
                          margin: const EdgeInsets.all(2),
                          decoration: BoxDecoration(
                            color: isSelected
                                ? AppColors.darkNavy
                                : isToday
                                ? AppColors.primaryGreen
                                : Colors.transparent,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Center(
                            child: Text(
                              '$day',
                              style: TextStyle(
                                color: isSelected
                                    ? Colors.white
                                    : AppColors.darkNavy,
                                fontWeight: isSelected || isToday
                                    ? FontWeight.bold
                                    : FontWeight.normal,
                                fontSize: 14,
                              ),
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  static String _monthName(int month) {
    const names = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ];
    return names[month - 1];
  }
}

class _ProgressCard extends StatelessWidget {
  final ProgressData data;
  const _ProgressCard({required this.data});

  @override
  Widget build(BuildContext context) {
    final progress = data.totalGoal > 0
        ? (data.totalGoal - data.remaining) / data.totalGoal
        : 0.0;

    return Container(
      decoration: const BoxDecoration(
        color: Color(0xFF9FBAF1),
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(12),
          topRight: Radius.circular(12),
          bottomLeft: Radius.circular(12),
          bottomRight: Radius.circular(12),
        ),
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 20, 16, 16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                _buildStat(data.eaten.toString(), 'Eaten'),
                _ThreeQuarterArc(
                  progress: progress.clamp(0.5, 1.0),
                  centerLabel: data.remaining.toString(),
                  subLabel: 'Remaining',
                ),
                _buildStat(data.burned.toString(), 'Burned'),
              ],
            ),
          ),
          Row(
            children: [
              _buildMacro(
                'Carbs',
                '${_fmt(data.carbsConsumed)}/${_fmt(data.carbsGoal)} g',
                const Color(0xFFEAF3B2),
                const BorderRadius.only(bottomLeft: Radius.circular(12)),
              ),
              _buildMacro(
                'Protein',
                '${_fmt(data.proteinConsumed)}/${_fmt(data.proteinGoal)} g',
                const Color(0xFFDCE4A8),
                BorderRadius.zero,
              ),
              _buildMacro(
                'Fat',
                '${_fmt(data.fatConsumed)}/${_fmt(data.fatGoal)} g',
                const Color(0xFFCAD297),
                const BorderRadius.only(bottomRight: Radius.circular(12)),
              ),
            ],
          ),
        ],
      ),
    );
  }

  String _fmt(double v) =>
      v == v.truncateToDouble() ? v.toInt().toString() : v.toStringAsFixed(1);

  Widget _buildStat(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: AppColors.darkNavy,
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: AppColors.darkNavy,
            fontSize: 12,
            fontWeight: FontWeight.w400,
          ),
        ),
      ],
    );
  }

  Widget _buildMacro(
    String label,
    String value,
    Color bgColor,
    BorderRadius borderRadius,
  ) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 4),
        decoration: BoxDecoration(color: bgColor, borderRadius: borderRadius),
        child: Column(
          children: [
            Text(
              label,
              style: const TextStyle(
                color: AppColors.darkNavy,
                fontSize: 14,
                fontWeight: FontWeight.w500,
              ),
            ),
            Text(
              value,
              style: const TextStyle(
                color: AppColors.darkNavy,
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ThreeQuarterArc extends StatelessWidget {
  final double progress;
  final String centerLabel;
  final String subLabel;

  const _ThreeQuarterArc({
    required this.progress,
    required this.centerLabel,
    required this.subLabel,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 130,
      height: 130,
      child: CustomPaint(
        painter: _ArcPainter(progress: progress),
        child: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                centerLabel,
                style: const TextStyle(
                  color: AppColors.darkNavy,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  height: 1.1,
                ),
              ),
              Text(
                subLabel,
                style: const TextStyle(
                  color: AppColors.darkNavy,
                  fontSize: 13,
                  fontWeight: FontWeight.w400,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ArcPainter extends CustomPainter {
  final double progress;
  const _ArcPainter({required this.progress});

  @override
  void paint(Canvas canvas, Size size) {
    const startAngle = 135.0 * math.pi / 180;
    const sweepFull = 270.0 * math.pi / 180;

    final rect = Rect.fromLTWH(12, 12, size.width - 24, size.height - 24);

    final trackPaint = Paint()
      ..color = AppColors.darkNavy.withValues(alpha: 0.12)
      ..strokeWidth = 14
      ..style = PaintingStyle.stroke
      ..strokeCap = StrokeCap.round;

    canvas.drawArc(rect, startAngle, sweepFull, false, trackPaint);

    if (progress > 0) {
      final progressPaint = Paint()
        ..color = AppColors.darkNavy.withValues(alpha: 0.40)
        ..strokeWidth = 14
        ..style = PaintingStyle.stroke
        ..strokeCap = StrokeCap.round;

      canvas.drawArc(
        rect,
        startAngle,
        sweepFull * progress,
        false,
        progressPaint,
      );
    }
  }

  @override
  bool shouldRepaint(_ArcPainter old) => old.progress != progress;
}

class _MealCardGroup extends StatelessWidget {
  final List<MealData> meals;
  const _MealCardGroup({required this.meals});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: AppColors.backgroundCream,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: meals.asMap().entries.map((e) {
          final idx = e.key;
          final meal = e.value;
          BorderRadius radius;
          if (idx == 0) {
            radius = const BorderRadius.only(
              topLeft: Radius.circular(15),
              topRight: Radius.circular(15),
            );
          } else if (idx == meals.length - 1) {
            radius = const BorderRadius.only(
              bottomLeft: Radius.circular(15),
              bottomRight: Radius.circular(15),
            );
          } else {
            radius = BorderRadius.zero;
          }
          final bool isLast = idx == meals.length - 1;

          return _MealCard(meal: meal, radius: radius);
        }).toList(),
      ),
    );
  }
}

class _MealCard extends StatelessWidget {
  final MealData meal;
  final BorderRadius radius;

  const _MealCard({required this.meal, required this.radius});

  @override
  Widget build(BuildContext context) {
    final isCream = meal.title == 'Breakfast';

    Color textColor = AppColors.darkNavy;
    Color iconColor = Colors.white;
    Color iconBgColor = AppColors.darkNavy;

    if (meal.title == 'Dinner') {
      textColor = const Color(0xFFEAF3B2);
      iconColor = const Color(0xFFEAF3B2);
    } else if (meal.title == 'Snack') {
      textColor = const Color(0xFFFDFDF5);
      iconColor = const Color(0xFFFDFDF5);
    }

    return GestureDetector(
      onTap: () {},
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
        decoration: BoxDecoration(color: meal.color, borderRadius: radius),
        child: Row(
          children: [
            CircleAvatar(
              backgroundColor: iconBgColor,
              radius: 18,
              child: Icon(meal.icon, color: iconColor, size: 20),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text(
                        meal.title,
                        style: TextStyle(
                          color: textColor,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                      const SizedBox(width: 4),
                      Icon(
                        Icons.chevron_right_rounded,
                        color: textColor,
                        size: 20,
                      ),
                    ],
                  ),
                  Text(
                    '${meal.consumed}/${meal.goal} kcal',
                    style: TextStyle(
                      color: textColor,
                      fontSize: 12,
                      fontWeight: FontWeight.w800,
                    ),
                  ),
                ],
              ),
            ),
            GestureDetector(
              onTap: () {},
              child: Container(
                padding: const EdgeInsets.all(5),
                decoration: BoxDecoration(
                  color: isCream
                      ? AppColors.darkNavy.withValues(alpha: 0.1)
                      : iconBgColor,
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.add,
                  color: isCream ? AppColors.darkNavy : iconColor,
                  size: 17,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
