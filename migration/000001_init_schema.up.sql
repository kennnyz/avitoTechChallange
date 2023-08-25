CREATE TABLE IF NOT EXISTS users (
    id int PRIMARY KEY
);

create table if not exists segments (
    segment_name varchar(255) PRIMARY KEY
);

create table if not exists user_segments (
                                             user_id int,
                                             segment_name varchar(255),
                                             PRIMARY KEY (user_id, segment_name),
                                             FOREIGN KEY (user_id) REFERENCES users(id),
                                             FOREIGN KEY (segment_name) REFERENCES segments(segment_name)
);

INSERT INTO users (id) VALUES (1000), (1002), (1003), (1004), (1005), (1006), (1007), (1008), (1009), (10100), (12342), (12343), (12344), (12345), (12346),
                              (12347), (12348), (12349), (123410), (123411), (123412), (123413), (123414), (123415), (123416), (123417), (123418), (123419),
                              (123420), (123421), (123422), (123423), (123424), (123425), (123426), (123427), (123428), (123429), (123430), (123431),
                              (123432), (123433), (123434), (123435), (123436), (123437), (123438), (123439), (123440), (123441), (123442), (123443), (123444),
                              (123445), (123446), (123447), (123448), (123449), (123450), (123451), (123452), (123453), (123454), (123455), (123456), (123457),
                              (123458), (123459), (123460), (123461), (123462), (123463), (123464), (123465), (123466), (123467), (123468), (123469), (123470),
                              (123471), (123472), (123473), (123474), (123475), (123476), (123477), (123478), (123479), (123480), (123481), (123482), (123483),
                              (123484), (123485), (123486), (123487), (123488), (123489), (123490), (123491), (123492), (123493), (123494), (123495), (123496),
                              (123497), (123498), (123499), (123500), (123501), (123502), (123503), (123504), (123505), (123506), (123507), (123508), (123509),
                              (123510), (123511), (123512), (123513), (123514), (123515), (123516), (123517), (123518), (123519), (123520), (123521), (123522),
                              (123523), (123524), (123525), (123526), (123527);

INSERT INTO segments (segment_name) VALUES ('AVITO_VOICE_MESSAGES'), ('AVITO_PERFORMANCE_VAS'), ('AVITO_DISCOUNT_30'), ('AVITO_DISCOUNT_50');

INSERT INTO user_segments (user_id, segment_name) VALUES (1000, 'AVITO_VOICE_MESSAGES'), (1000, 'AVITO_DISCOUNT_30'),
                                                          (1000, 'AVITO_DISCOUNT_50'), (1002, 'AVITO_DISCOUNT_30'), (1002, 'AVITO_DISCOUNT_50'),
                                                          (1003, 'AVITO_VOICE_MESSAGES'),
                                                          (1003, 'AVITO_PERFORMANCE_VAS'), (1003, 'AVITO_DISCOUNT_30'), (1003, 'AVITO_DISCOUNT_50'),
                                                          (1004, 'AVITO_VOICE_MESSAGES'), (1004, 'AVITO_PERFORMANCE_VAS'), (1004, 'AVITO_DISCOUNT_30'),
                                                          (1004, 'AVITO_DISCOUNT_50'), (1005, 'AVITO_VOICE_MESSAGES'),
                                                          (1005, 'AVITO_DISCOUNT_30'), (1005, 'AVITO_DISCOUNT_50'), (1006, 'AVITO_VOICE_MESSAGES'),
                                                          (1006, 'AVITO_PERFORMANCE_VAS'), (1006, 'AVITO_DISCOUNT_30'), (1006, 'AVITO_DISCOUNT_50'),
                                                          (1007, 'AVITO_PERFORMANCE_VAS'), (1007, 'AVITO_DISCOUNT_30'),
                                                          (1007, 'AVITO_DISCOUNT_50'), (1008, 'AVITO_VOICE_MESSAGES'), (1008, 'AVITO_PERFORMANCE_VAS'),
                                                          (1008, 'AVITO_DISCOUNT_30'),  (1009, 'AVITO_VOICE_MESSAGES'),
                                                          (1009, 'AVITO_PERFORMANCE_VAS'), (1009, 'AVITO_DISCOUNT_30'), (1009, 'AVITO_DISCOUNT_50'),
                                                          (10100, 'AVITO_VOICE_MESSAGES'),  (10100, 'AVITO_DISCOUNT_30'),
                                                          (10100, 'AVITO_DISCOUNT_50'), (12342, 'AVITO_VOICE_MESSAGES');