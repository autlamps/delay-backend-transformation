# Python Implementation of Transformation
Backend service responsible for inserting static data into the database. This branch is a python implementation. It was written as quickly as possible to get a minimum viable product to allow the 
collection service to be made and to test the golang version against. It has no tests and makes no guarantees of correctness.

__Warning__: this is seriously slow!

## Usage

1. Change psql access url to suit needs
2. Apply schema on db
3. Place next to a folder called gtfs containing gtfs static data in a txt file
4. Run with ``python3 transform.py``
5. Success?

## License

Copyright (C) 2017 Dharyin Colbert, Izaac Crooke, Dominic Porter, Hayden Woodhead

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
