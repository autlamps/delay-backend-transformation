import csv
import uuid
import psycopg2

all_agency = dict()
all_route = dict()
all_calendar = dict()
all_trip = dict()
all_stop = dict()


class Agency:

    def __init__(self, agency_id, gtfs_agency_id, agency_name):
        self.agency_id = str(agency_id)
        self.gtfs_agency_id = gtfs_agency_id
        self.agency_name = agency_name


class Route:

    def __init__(self, route_id, gtfs_route_id, agency_id, route_short_name, route_long_name):
        self.route_id = str(route_id)
        self.gtfs_route_id = gtfs_route_id
        self.agency_id = agency_id
        self.route_short_name = route_short_name
        self.route_long_name = route_long_name


class Trip:

    def __init__(self, trip_id, gtfs_trip_id, route_id, service_id, trip_headsign):
        self.trip_id = str(trip_id)
        self.gtfs_trip_id = gtfs_trip_id
        self.route_id = route_id
        self.service_id = service_id
        self.trip_headsign = trip_headsign


class Calendar:

    def __init__(self, service_id, gtfs_service_id, start_date, end_date, mon, tue, wed, thu, fri, sat, sun):
        self.service_id = str(service_id)
        self.gtfs_service_id = gtfs_service_id
        self.start_date = start_date
        self.end_date = end_date
        self.mon = str(mon)
        self.tue = str(tue)
        self.wed = str(wed)
        self.thu = str(thu)
        self.fri = str(fri)
        self.sat = str(sat)
        self.sun = str(sun)


class Stop:

    def __init__(self, stop_id, stop_name, stop_lat, stop_long):
        self.stop_id = str(stop_id)
        self.stop_name = stop_name
        self.stop_lat = stop_lat
        self.stop_long = stop_long


class StopTime:

    def __init__(self, stoptime_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence):
        self.stoptime_id = str(stoptime_id)
        self.trip_id = trip_id
        self.arrival_time = arrival_time
        self.departure_time = departure_time
        self.stop_id = stop_id
        self.stop_sequence = stop_sequence


def read_agency(cur):
    with open('gtfs/agency.txt', newline='') as agency_file:
        next(agency_file, None)  # Skip first line as this has column names

        reader = csv.reader(agency_file)

        for row in reader:
            ag = Agency(uuid.uuid4(), row[2], row[3])
            all_agency[ag.gtfs_agency_id] = ag
            cur.execute("insert into agency (agency_id, gtfs_agency_id, agency_name) VALUES (%s, %s, %s)",
                        (ag.agency_id, ag.gtfs_agency_id, ag.agency_name))


def read_routes(cur):
    with open('gtfs/routes.txt', newline='') as routes_file:
        next(routes_file, None)

        reader = csv.reader(routes_file)

        for row in reader:
            r = Route(uuid.uuid4(), row[4], all_agency[row[3]].agency_id, row[6], row[0])
            all_route[r.gtfs_route_id] = r

            cur.execute("insert into routes (route_id, gtfs_route_id, agency_id, route_short_name, route_long_name) "
                        "VALUES (%s, %s, %s, %s, %s)", (r.route_id, r.gtfs_route_id, r.agency_id, r.route_short_name,
                                                        r.route_long_name))


def read_calendar(cur):
    with open('gtfs/calendar.txt', newline='') as cal_file:
        next(cal_file, None)

        reader = csv.reader(cal_file)

        for row in reader:
            c = Calendar(uuid.uuid4(), row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9])
            all_calendar[c.gtfs_service_id] = c

            cur.execute("insert into calendar (service_id, gtfs_service_id, start_date, end_date, monday, tuesday,"
                        " wednesday, thursday, friday, saturday, sunday) "
                        "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
                        (c.service_id, c.gtfs_service_id, c.start_date, c.end_date, c.mon, c.tue, c.wed, c.thu,
                         c.fri, c.sat, c.sun))


def read_trips(cur):
    with open('gtfs/trips.txt', newline='') as trips_file:
        next(trips_file, None)

        reader = csv.reader(trips_file)

        for row in reader:
            t = Trip(uuid.uuid4(), row[6], all_route[row[1]].route_id, all_calendar[row[5]].service_id, row[3])
            all_trip[t.gtfs_trip_id] = t

            cur.execute("insert into trips (trip_id, route_id, service_id, gtfs_trip_id, trip_headsign)"
                        " VALUES (%s, %s, %s, %s, %s)",
                        (t.trip_id, t.route_id, t.service_id, t.gtfs_trip_id, t.trip_headsign))


def read_stops(cur):
    with open('gtfs/stops.txt', newline='') as stops_file:
        next(stops_file, None)

        reader = csv.reader(stops_file)

        for row in reader:
            s = Stop(uuid.uuid4(), stop_name=row[6], stop_lat=row[0], stop_long=row[2])
            all_stop[row[3]] = s

            cur.execute("insert into stops (stop_id, stop_name, stop_lat, stop_lon) VALUES (%s, %s, %s, %s)",
                        (s.stop_id, s.stop_name, s.stop_lat, s.stop_long))


def read_stoptimes(cur):
    with open('gtfs/stop_times.txt', newline='') as stoptime_file:
        next(stoptime_file, None)

        reader = csv.reader(stoptime_file)

        for row in reader:
            st = StopTime(uuid.uuid4(), trip_id=all_trip[row[0]].trip_id, arrival_time=row[1], departure_time=row[2],
                          stop_id=all_stop[row[3]].stop_id, stop_sequence=row[4])

            # Pull the hour part from arrival and departure time
            at = int(st.arrival_time[:2])
            dt = int(st.departure_time[:2])

            if at >= 24:  # If over 24
                at = at - 24   # remove 24
                # insert it back into the hour section of the time string
                st.arrival_time = "{}:{}".format(str(at), st.arrival_time[3:])

            if dt >= 24:
                dt = dt - 24
                st.departure_time = "{}:{}".format(str(dt), st.departure_time[3:])

            cur.execute("insert into stop_times (stoptime_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence) "
                        "VALUES (%s, %s, %s, %s, %s, %s)", (st.stoptime_id, st.trip_id, st.arrival_time,
                                                            st.departure_time, st.stop_id, st.stop_sequence))


def main():
    conn = psycopg2.connect("postgresql://postgres:mysecretpassword@172.17.0.2/postgres")
    cur = conn.cursor()

    read_agency(cur)
    read_routes(cur)
    read_calendar(cur)
    read_trips(cur)
    read_stops(cur)
    read_stoptimes(cur)

    conn.commit()
    cur.close()
    conn.close()

if __name__ == '__main__':
    main()