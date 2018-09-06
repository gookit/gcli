package progress

import "time"

func ExampleLoadingSpinner() {
	s := LoadingSpinner(RandomCharsTheme(), 100*time.Millisecond)

	s.Start("%s work handling ... ...")
	// Run for some time to simulate work
	time.Sleep(4 * time.Second)
	s.Stop("work handle complete")
}

func ExampleRoundTripSpinner() {
	s := RoundTripSpinner(RandomCharTheme(), 100*time.Millisecond)

	s.Start("%s work handling ... ...")
	// Run for some time to simulate work
	time.Sleep(4 * time.Second)
	s.Stop("work handle complete")
}
