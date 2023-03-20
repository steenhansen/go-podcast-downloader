package processes

import (
	"context"
)

/*
func keyboardFailure(err error) models.PodcastResults {
	locStat := errors.New("DOWNLOADMEDIA()")
	multiErr := errors.Join(locStat, err)
	emptyPodcastResults := misc.EmptyPodcastResults(false, multiErr)
	return emptyPodcastResults

}
*/

func firstErr(err error, seriousStream <-chan error) error {
	for i := 0; i < len(seriousStream); i++ {
		err := <-seriousStream
		if err != nil {
			return err
		}
	}
	if err == context.Canceled {
		return nil
	}
	return err // like weird keyboard error, or disk read only?
}
