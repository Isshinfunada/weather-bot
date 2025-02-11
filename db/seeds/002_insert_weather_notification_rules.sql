-- +goose Up

INSERT INTO weather_notification_rules (weather_code, weather_description, is_notify_trigger) VALUES
('100', '晴', FALSE),
('101', '晴時々曇', FALSE),
('110', '晴後時々曇', FALSE),
('111', '晴後曇', FALSE),
('130', '朝の内霧後晴', FALSE),
('131', '晴明け方霧', FALSE),
('132', '晴朝夕曇', FALSE),
('200', '曇', FALSE),
('201', '曇時々晴', FALSE),
('209', '霧', FALSE),
('210', '曇後時々晴', FALSE),
('211', '曇後晴', FALSE),
('223', '曇日中時々晴', FALSE),
/*雨を含むパターン*/
('102', '晴一時雨', TRUE),
('103', '晴時々雨', TRUE),
('104', '晴一時雪', TRUE),
('105', '晴時々雪', TRUE),
('106', '晴一時雨か雪', TRUE),
('107', '晴時々雨か雪', TRUE),
('108', '晴一時雨か雷雨', TRUE),
('112', '晴後一時雨', TRUE),
('113', '晴後時々雨', TRUE),
('114', '晴後雨', TRUE),
('115', '晴後一時雪', TRUE),
('116', '晴後時々雪', TRUE),
('117', '晴後雪', TRUE),
('118', '晴後雨か雪', TRUE),
('119', '晴後雨か雷雨', TRUE),
('120', '晴朝夕一時雨', TRUE),
('121', '晴朝の内一時雨', TRUE),
('122', '晴夕方一時雨', TRUE),
('123', '晴山沿い雷雨', TRUE),
('124', '晴山沿い雪', TRUE),
('125', '晴午後は雷雨', TRUE),
('126', '晴昼頃から雨', TRUE),
('127', '晴夕方から雨', TRUE),
('128', '晴夜は雨', TRUE),
('140', '晴時々雨で雷を伴う', TRUE),
('160', '晴一時雪か雨', TRUE),
('170', '晴時々雪か雨', TRUE),
('181', '晴後雪か雨', TRUE),
('202', '曇一時雨', TRUE),
('203', '曇時々雨', TRUE),
('204', '曇一時雪', TRUE),
('205', '曇時々雪', TRUE),
('206', '曇一時雨か雪', TRUE),
('207', '曇時々雨か雪', TRUE),
('208', '曇一時雨か雷雨', TRUE),
('212', '曇後一時雨', TRUE),
('213', '曇後時々雨', TRUE),
('214', '曇後雨', TRUE),
('215', '曇後一時雪', TRUE),
('216', '曇後時々雪', TRUE),
('217', '曇後雪', TRUE),
('218', '曇後雨か雪', TRUE),
('219', '曇後雨か雷雨', TRUE),
('220', '曇朝夕一時雨', TRUE),
('221', '曇朝の内一時雨', TRUE),
('222', '曇夕方一時雨', TRUE),
('224', '曇昼頃から雨', TRUE),
('225', '曇夕方から雨', TRUE),
('226', '曇夜は雨', TRUE),
('228', '曇昼頃から雪', TRUE),
('229', '曇夕方から雪', TRUE),
('230', '曇夜は雪', TRUE),
('231', '曇海上海岸は霧か霧雨', TRUE),
('240', '曇時々雨で雷を伴う', TRUE),
('250', '曇時々雪で雷を伴う', TRUE),
('260', '曇一時雪か雨', TRUE),
('270', '曇時々雪か雨', TRUE),
('281', '曇後雪か雨', TRUE),
('300', '雨', TRUE),
('301', '雨時々晴', TRUE),
('302', '雨時々止む', TRUE),
('303', '雨時々雪', TRUE),
('304', '雨か雪', TRUE),
('306', '大雨', TRUE),
('308', '雨で暴風を伴う', TRUE),
('309', '雨一時雪', TRUE),
('311', '雨後晴', TRUE),
('313', '雨後曇', TRUE),
('314', '雨後時々雪', TRUE),
('315', '雨後雪', TRUE),
('316', '雨か雪後晴', TRUE),
('317', '雨か雪後曇', TRUE),
('320', '朝の内雨後晴', TRUE),
('321', '朝の内雨後曇', TRUE),
('322', '雨朝晩一時雪', TRUE),
('323', '雨昼頃から晴', TRUE),
('324', '雨夕方から晴', TRUE),
('325', '雨夜は晴', TRUE),
('326', '雨夕方から雪', TRUE),
('327', '雨夜は雪', TRUE),
('328', '雨一時強く降る', TRUE),
('329', '雨一時みぞれ', TRUE),
('340', '雪か雨', TRUE),
('350', '雨で雷を伴う', TRUE),
('361', '雪か雨後晴', TRUE),
('371', '雪か雨後曇', TRUE),
('400', '雪', TRUE),
('401', '雪時々晴', TRUE),
('402', '雪時々止む', TRUE),
('403', '雪時々雨', TRUE),
('405', '大雪', TRUE),
('406', '風雪強い', TRUE),
('407', '暴風雪', TRUE),
('409', '雪一時雨', TRUE),
('411', '雪後晴', TRUE),
('413', '雪後曇', TRUE),
('414', '雪後雨', TRUE),
('420', '朝の内雪後晴', TRUE),
('421', '朝の内雪後曇', TRUE),
('422', '雪昼頃から雨', TRUE),
('423', '雪夕方から雨', TRUE),
('425', '雪一時強く降る', TRUE),
('426', '雪後みぞれ', TRUE),
('427', '雪一時みぞれ', TRUE),
('450', '雪で雷を伴う', TRUE);
