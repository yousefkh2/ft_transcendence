<template>
	<LegalLayout title="Privacy Policy" updated="October 26, 2026">
		<p>
			This policy explains which personal data LingoDash collects when you use it, why we need it,
			who it is shared with and what rights you have. We only collect what the game needs to work.
		</p>

		<h2>1. Who we are</h2>
		<p>
			LingoDash is a multiplayer language-learning game built by a team of students as part of the
			ft_transcendence project of the 42 curriculum at 42 Heilbronn, Germany. It is a non-commercial
			student project. The project team is responsible for the processing of your data described here.
		</p>
		<p>
			You can reach the team through the project repository:
			<a :href="repoUrl" target="_blank" rel="noopener noreferrer">github.com/yousefkh2/ft_transcendence</a>.
		</p>

		<h2>2. Data we collect</h2>
		<h3>Account data</h3>
		<ul>
			<li>Your username and email address.</li>
			<li>
				Your password, which we never store in readable form. It is saved only as a salted bcrypt hash, so
				nobody on the team can see it.
			</li>
			<li>The date your account was created.</li>
		</ul>
		<h3>Game and progress data</h3>
		<ul>
			<li>Your XP.</li>
			<li>The rounds you play: game mode, language, start and end time, your role and the result.</li>
			<li>
				Corrections from your rounds that can be saved as flashcards: what you said, the corrected phrase
				and a short explanation.
			</li>
		</ul>
		<h3>Friends</h3>
		<ul>
			<li>Friend requests you send or receive, their status and the date they were made.</li>
		</ul>
		<h3>Data stored in your browser</h3>
		<ul>
			<li>
				After you log in, a login token is saved in your browser's local storage so you stay logged in. It
				expires after 24 hours and is removed when you log out.
			</li>
		</ul>
		<p>
			LingoDash does <strong>not</strong> use cookies, analytics, advertising or any other tracking tools.
		</p>

		<h2>3. Your voice and speech recognition</h2>
		<p>
			Speaking is the core of the game. When you enable your microphone during a round, your voice is sent
			directly from your browser to the Realtime API of OpenAI, which turns it into text. The resulting
			transcript is passed through our server to the other player in your room, so the game can react to
			what you said.
		</p>
		<ul>
			<li>We do <strong>not</strong> record or store your audio.</li>
			<li>We do <strong>not</strong> store the transcripts of your rounds on our servers.</li>
			<li>
				OpenAI processes the audio for us as a service provider. Its handling of data is described in the
				<a href="https://openai.com/policies/privacy-policy/" target="_blank" rel="noopener noreferrer">OpenAI Privacy Policy</a>.
			</li>
		</ul>
		<p>You decide when your microphone is on. Your browser always asks for permission first.</p>

		<h2>4. Why we use your data</h2>
		<ul>
			<li>To create your account and log you in.</li>
			<li>To run game rounds, match you with other players and keep track of your XP and rounds.</li>
			<li>To let you add friends and see your friends list.</li>
			<li>To turn your speech into text during a round.</li>
			<li>To keep the service secure, for example by protecting accounts against misuse.</li>
		</ul>
		<p>
			Under the EU General Data Protection Regulation (GDPR), we process this data because it is needed to
			provide the service you signed up for (Art. 6(1)(b) GDPR) and, for security measures, because of our
			legitimate interest in a safe service (Art. 6(1)(f) GDPR).
		</p>

		<h2>5. Who can see your data</h2>
		<ul>
			<li>Other players can see your username, for example in rounds and friend requests.</li>
			<li>
				Logged-in players can open your public profile. It shows your username, XP, the date you joined and
				your game stats (rounds played, wins and losses). Your email and your detailed match history stay
				private.
			</li>
			<li>Your friends can see that you are friends with them.</li>
			<li>Your email address and password hash are never shown to other players.</li>
			<li>The project team can access the database to run and maintain LingoDash.</li>
			<li>OpenAI receives your voice audio while your microphone is on, as described in section 3.</li>
		</ul>
		<p>We never sell your data or share it for advertising.</p>

		<h2>6. Storage and security</h2>
		<p>
			Your data is stored in a PostgreSQL database run by the project team. Passwords are hashed with
			bcrypt, login tokens expire after 24 hours, and the connection between your browser and LingoDash is
			encrypted with HTTPS.
		</p>

		<h2>7. How long we keep your data</h2>
		<p>
			We keep your data for as long as your account exists. When your account is deleted, your profile,
			friendships, rounds and flashcards are deleted with it.
		</p>

		<h2>8. Your rights</h2>
		<p>Under the GDPR you have the right to:</p>
		<ul>
			<li>get a copy of the data we hold about you,</li>
			<li>have incorrect data corrected,</li>
			<li>have your account and data deleted,</li>
			<li>restrict or object to certain processing,</li>
			<li>receive your data in a portable format.</li>
		</ul>
		<p>
			To use any of these rights, contact the team through the project repository linked in section 1. You
			can also file a complaint with a data protection authority, for example the State Commissioner for
			Data Protection and Freedom of Information Baden-Württemberg.
		</p>

		<h2>9. Age requirement</h2>
		<p>You must be at least 16 years old to create a LingoDash account.</p>

		<h2>10. Changes to this policy</h2>
		<p>
			If LingoDash starts collecting new data or uses it in a new way, we will update this page and the date
			at the top.
		</p>

		<p>
			See also our <RouterLink to="/terms">Terms of Service</RouterLink>.
		</p>
	</LegalLayout>
</template>

<script setup lang="ts">
import LegalLayout from './LegalLayout.vue';

const repoUrl = 'https://github.com/yousefkh2/ft_transcendence';
</script>
