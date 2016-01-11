Reyong Composer
===============

This is a simple composition tool that generates reyong patterns for the Balinese [beleganjur](https://en.wikipedia.org/wiki/Beleganjur) orchestra.

In beleganjur, two players play two reyong notes each. The lower two notes are called polos, and the upper two are called sangsih. A music notation system is used in which numbers signify notes, and dots signify rests. Here is an example of the type of music which this program is designed to create. It is a 32-beat repeating pattern:

    polos:   [1.212.12.12.1.21.12.121.2.12.1.2]
    sangsih: [43.4.34.34.343.434.34.43.34.343.]

Stylistic Constraints
---------------------

* No note or rest can repeat.
* There can be no more than three notes without a rest.
* Polos must start with a note, not a rest.
* No repeated note and rest pairs, i.e. note-rest-note-rest-note.
* Sangsih and polos cannot have a rest at the same time.
* 1 and 3 cannot be combined, and neither can 2 and 4.

TODO
----------

* Generate audio for the patterns.
