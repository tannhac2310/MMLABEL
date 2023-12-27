 -- số lượng in giấy trong ngày
 ALTER TABLE device_working_history
     ADD number_of_prints_per_day INT DEFAULT 0;
ALTER TABLE device_working_history
    ADD printing_time_per_day INT DEFAULT 0;
