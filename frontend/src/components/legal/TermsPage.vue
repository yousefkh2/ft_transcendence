<template>
	<LegalLayout title="Terms of Service" updated="October 26, 2026">
		<p>
			These terms are the rules for using LingoDash. By creating an account or playing a round, you agree to
			them. Please read them together with our <RouterLink to="/privacy">Privacy Policy</RouterLink>.
		</p>

		<h2>1. About LingoDash</h2>
		<p>
			LingoDash is a free multiplayer game for practicing languages by speaking with other learners. It is a
			non-commercial student project, built by a team of students as part of the ft_transcendence project
			of the 42 curriculum at 42 Heilbronn, Germany.
		</p>

		<h2>2. Who can use LingoDash</h2>
		<ul>
			<li>You must be at least 16 years old.</li>
			<li>You may only create one account, and only for yourself.</li>
			<li>The information you give us, such as your email address, has to be correct.</li>
		</ul>

		<h2>3. Your account</h2>
		<ul>
			<li>Keep your password secret. You are responsible for everything that happens with your account.</li>
			<li>Choose a username that does not insult anyone or pretend to be someone else.</li>
			<li>If you think someone else has access to your account, change your password and contact us.</li>
		</ul>

		<h2>4. Playing fair and being respectful</h2>
		<p>
			LingoDash only works if people feel safe speaking a language they are still learning. When you play,
			you agree not to:
		</p>
		<ul>
			<li>insult, harass, threaten or discriminate against other players,</li>
			<li>say or send sexual, hateful or violent content,</li>
			<li>share other people's personal information,</li>
			<li>cheat, use bots or scripts, or exploit bugs to gain XP or win rounds,</li>
			<li>attack, overload or try to break into LingoDash or other players' accounts,</li>
			<li>use LingoDash for advertising or anything illegal.</li>
		</ul>
		<p>Mistakes in the language are expected. Being unkind about them is not.</p>

		<h2>5. Your voice and what you say</h2>
		<ul>
			<li>
				During a round, your speech is turned into text and shared with the other player in your room, as
				described in the <RouterLink to="/privacy">Privacy Policy</RouterLink>.
			</li>
			<li>You are responsible for what you say in a round.</li>
			<li>
				What you say stays yours. You only allow us to process it as far as needed to run the game, for
				example to show it to your partner or to create flashcards for you.
			</li>
		</ul>

		<h2>6. Corrections and learning content</h2>
		<p>
			Feedback and corrections in LingoDash can be generated automatically and may contain mistakes. Treat
			them as a learning aid, not as an official language reference.
		</p>

		<h2>7. Availability</h2>
		<p>
			LingoDash is a student project and is provided "as is". Features can change, the service can be
			unavailable at times, and data such as XP may be reset during development. We do our best to keep it
			running, but we cannot guarantee it.
		</p>

		<h2>8. Suspending or closing accounts</h2>
		<ul>
			<li>
				If you break these terms, we may warn you, suspend your account or delete it, depending on how
				serious the problem is.
			</li>
			<li>
				You can stop using LingoDash at any time and delete your account yourself on your profile page.
			</li>
		</ul>

		<h2>9. Liability</h2>
		<p>
			LingoDash is free and provided without any warranty. We are not liable for damage resulting from its
			use, except where the law does not allow this to be excluded, in particular for damage caused
			intentionally or by gross negligence, or for injury to life, body or health.
		</p>

		<h2>10. Changes to these terms</h2>
		<p>
			We may update these terms when LingoDash changes. The date at the top shows the latest version. If you
			keep using LingoDash after a change, the new terms apply.
		</p>

		<h2>11. Contact and applicable law</h2>
		<p>
			You can reach the team through the project repository:
			<a :href="repoUrl" target="_blank" rel="noopener noreferrer">github.com/yousefkh2/ft_transcendence</a>.
			These terms are governed by the law of the Federal Republic of Germany.
		</p>
	</LegalLayout>
</template>

<script setup lang="ts">
import LegalLayout from './LegalLayout.vue';

const repoUrl = 'https://github.com/yousefkh2/ft_transcendence';
</script>
