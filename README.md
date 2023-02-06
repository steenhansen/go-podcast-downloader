



TESTS
  Test everything
    go test ./... -count=1 | grep -E "^ok.*s$" 

  Tests using real internet 
    go test ./src/internet-tests/... -count=1 | grep -E "^ok.*s$"


  Test using mocked internet
    go test ./... -run "^Test_.*$" -count=1 | grep -E "^ok.*s$" 

OPTIONAL ARGUMENTS
  --networkLoad sets the amount of network traffic
      go run ./ --networkLoad=high
      go run ./ --networkLoad=medium
      go run ./ --networkLoad=low 

  --fileLimit sets the maximum number of files to download from a podcast
      go run ./ --fileLimit=3


https://podcasts.apple.com/ca/podcast/the-history-of-the-twentieth-century/id1039714402

https://allthingscomedy.com/podcast/the-dollop

https://play.acast.com/s/the-rest-is-history-podcast

https://timesuckpodcast.com/Episodes

https://www.bbc.co.uk/programmes/articles/1xNChQmb9MRxSVcpChLttvr/witness-history-podcasts














