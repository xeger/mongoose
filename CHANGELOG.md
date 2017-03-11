v1
==

v0
==

Trying to keep pace with gomuti for now, and to support all of gomuti's features (mocks, spies, etc).

Known issues:

Mongoose panics in a less-than-helpful way if a call does not match. Should print the closest match in a way similar to gomuti's SpyMatcher.

Mongoose generates panicky code if a mock returns fewer values than expected, or if results are not the expected type. Should check for this
and panic with a more informative, user-friendly message. Maybe try to auto-convert too, using reflection ... hmm....
