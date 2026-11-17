import type { AchievementId } from '../api/progress';
import type { IconName } from '../components/base/BaseIcon.vue';

export const ACHIEVEMENT_ICONS: Record<AchievementId, IconName> = {
	first_round: 'play',
	first_win: 'trophy',
	first_loss: 'refresh',
	quick_solve: 'clock',
	lightning: 'zap',
	flawless: 'target',
	streak_3: 'flame',
	rounds_10: 'award',
	wins_10: 'trophy',
	level_5: 'star',
	level_10: 'star',
	polyglot: 'globe',
	team_player: 'users',
};
