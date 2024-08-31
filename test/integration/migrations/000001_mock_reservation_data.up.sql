INSERT INTO reservations (id, room_id, start_time, end_time, created_at, updated_at) VALUES
-- 1st room, 3 resevations: 09:00 <===> 10:00 <===> 12:00 ..(free 2 hour slot).. 14:00 <===> 16:00 
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5e', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f', '2024-08-31 14:00:00', '2024-08-31 16:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b60', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f', '2024-08-31 10:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b62', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f', '2024-08-31 09:00:00', '2024-08-31 10:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
-- 2nd room, 3 resevations: 08:00 <===> 10:00 <===> 13:00 <===> 15:00 
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b64', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b65', '2024-08-31 13:00:00', '2024-08-31 15:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b66', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b65', '2024-08-31 08:00:00', '2024-08-31 10:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b68', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b65', '2024-08-31 10:00:00', '2024-08-31 13:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
-- other reservations per room:
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6a', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6b', '2024-08-31 06:00:00', '2024-08-31 08:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6c', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6d', '2024-08-31 05:00:00', '2024-08-31 07:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6e', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6f', '2024-08-31 04:00:00', '2024-08-31 06:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b70', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b71', '2024-08-31 03:00:00', '2024-08-31 05:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b72', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b73', '2024-08-31 02:00:00', '2024-08-31 04:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b74', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b75', '2024-08-31 01:00:00', '2024-08-31 03:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00'),
('018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b76', '018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b77', '2024-08-31 00:00:00', '2024-08-31 02:00:00', '2024-08-31 12:00:00', '2024-08-31 12:00:00');