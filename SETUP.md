# Setup: run LingoDash and play with 2 PCs

This guide gets the project running on your computer and shows how a second
computer can join, so you can test the game (and the voice chat) with two players.

One computer is the **host**: it runs everything. The second computer only needs
a web browser.

---

## 1. What you need (host computer only)

- **Git**
- **Docker Desktop** – start it and leave it running
- **make** and **openssl** – already installed on macOS and most Linux systems
- **Godot 4** – only if you want to change the game

The second computer needs nothing except a browser (Chrome or Firefox) and a
microphone.

---

## 2. Get the project

```sh
git clone https://github.com/yousefkh2/ft_transcendence.git
cd ft_transcendence
```

---

## 3. Create your settings file

```sh
cp .env.example .env
```

Open `.env` and change one line:

```sh
JWT_SECRET=pick-any-long-random-text-here
```

Leave the rest as it is. Keep `VITE_API_URL=` **empty**, otherwise the second
PC will not work.

---

## 4. Start everything

```sh
make up
```

The first start takes a few minutes. When it is done it prints something like:

```text
Open https://localhost on this computer.
Other computers on the network: https://192.168.1.20
The certificate is self-signed: accept the browser warning once on each computer.
```

Write down the second address (the one with numbers). That is your **host address**.

You can print it again any time with:

```sh
make url
```

---

## 5. Open the website on the host

1. Open **https://localhost** in your browser.
2. The browser shows a security warning (the certificate is made by us, not by a
   real certificate company). Click **Advanced → Proceed / Accept the risk**.
3. Register an account, log in, and click **Play**.

---

## 6. Join from the second PC

1. Make sure both computers are on the **same Wi-Fi / network**.
2. On the second PC, open the **host address** from step 4, for example
   **https://192.168.1.20**.
3. Accept the security warning again (same as step 5).
4. Register a **different** account and log in.

---

## 7. Play together

1. **Player 1:** in the game, create a lobby. You get a lobby code.
2. **Player 2:** enter that code and join the lobby.
3. Both browsers ask for the **microphone** – click **Allow**.
4. At the top of the page you see the voice status. Wait until it says
   **Voice connected** – now you can talk to each other.
5. Each player picks a role (**Mission Control** or **On Site**) and presses
   **Ready**. You can only get ready when your microphone is connected.
6. When both are ready, the match starts.

---

## Something does not work?

| Problem | What to do |
|---|---|
| `make up` fails with `JWT_SECRET must be set` | You skipped step 3. Set `JWT_SECRET` in `.env`. |
| `make up` fails, "Cannot connect to the Docker daemon" | Start Docker Desktop and try again. |
| The second PC cannot open the host address | Both PCs must be on the same network. Check that the host's firewall allows incoming connections for Docker. Run `make url` again – the IP can change when you switch networks. |
| The page loads but login / the game does not work on the second PC | Check that `VITE_API_URL=` is empty in `.env`, then run `make restart`. |
| "No microphone" | Click the lock / site-settings icon next to the address bar, allow the microphone, then leave the lobby and join it again. Always use **https://**, not http://. |
| Stuck on "Waiting for voice..." | The website works but the two PCs cannot reach each other directly. School and company networks often block this. Try a phone hotspot or a home network. |
| You hear yourself twice / echo | Use headphones, especially when both players are in the same room. |

---

## For developers

**Changed the Godot game?** Export it again from Godot (Project → Export →
Web). It is saved to `frontend/public/`, and the website picks it up right away.

**Testing from the Godot editor** (without the website): the editor cannot join
the voice chat, so readying up is blocked. Put this in `.env` and run
`make restart`:

```sh
VOICE_REQUIRED=false
```

**Useful commands:**

```sh
make up        # build and start everything
make down      # stop everything
make restart   # stop and start again (after changing .env)
make logs      # see what the services are doing
make url       # show the addresses to open
make test      # run the backend tests
```
