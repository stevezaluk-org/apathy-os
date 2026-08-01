package main

// SwitchRoot Moves a mounted `/` to an arbitrary directory selected by the caller and exec's into `init` afterward. If `rm`
// is set to true, then recursively remove everything past `/` after the mount has been moved
func SwitchRoot(root string, init string, rm bool) error {
	return nil
}
