ALTER TABLE `room_allotment` CHANGE COLUMN `from_date` `start_date` DATE NOT NULL;
ALTER TABLE `room_allotment` CHANGE COLUMN `to_date` `end_date` DATE NOT NULL;