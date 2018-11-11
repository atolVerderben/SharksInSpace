package spaceshark

func finalizeAudio() error {
	return nil
}

func loadAudio() error {
	//soundtrack
	AudioPlayer.AddSongFromFile("map", "assets/audio/music/map.wav")
	AudioPlayer.AddSongFromFile("station", "assets/audio/music/station.wav")
	AudioPlayer.AddSongFromFile("thruster", "assets/audio/thruster.wav")

	//sound effects
	AudioPlayer.AddSoundEffectFromFile("laserHit", "assets/audio/laser-hit.wav", 0.2)
	AudioPlayer.AddSoundEffectFromFile("playerDie", "assets/audio/player-die-explosion.wav", 0.2)
	AudioPlayer.AddSoundEffectFromFile("playerHit", "assets/audio/player-hit.wav", 0.2)
	AudioPlayer.AddSoundEffectFromFile("sharkExplode", "assets/audio/shark-explode.wav", 0.2)
	AudioPlayer.AddSoundEffectFromFile("shootLaser", "assets/audio/shoot-laser.wav", 0.2)

	return nil
}
