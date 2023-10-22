CREATE USER 'tca-students' IDENTIFIED BY 'strongpassword';

create table if not exists actions
(
    action      int auto_increment primary key,
    name        varchar(50) not null,
    description mediumtext  not null,
    color       varchar(6)  not null
) collate = utf8mb4_unicode_ci
  auto_increment = 19;

create table if not exists alarm_ban
(
    ban     int auto_increment primary key,
    created timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    ip      binary(16)                                              not null,
    constraint ip unique (ip)
) collate = utf8mb4_unicode_ci;

create table if not exists alarm_log
(
    alarm    int auto_increment primary key,
    created  timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    message  text                                                    not null,
    send     int                                                     not null,
    received int                                                     not null,
    test     enum ('true', 'false')      default 'false'             not null,
    ip       binary(16)                                              not null
) collate = utf8mb4_unicode_ci;

create table if not exists barrierFree_moreInfo
(
    id       int(11) unsigned auto_increment primary key,
    title    varchar(32)  null,
    category varchar(11)  null,
    url      varchar(128) null
) charset = utf8
  auto_increment = 11;

create table if not exists barrierFree_persons
(
    id         int(11) unsigned auto_increment primary key,
    name       varchar(40) null,
    telephone  varchar(32) null,
    email      varchar(32) null,
    faculty    varchar(32) null,
    office     varchar(16) null,
    officeHour varchar(16) null,
    tumID      varchar(24) null
) charset = utf8
  auto_increment = 19;

create table if not exists card_type
(
    card_type int auto_increment primary key,
    title     varchar(255) null
) auto_increment = 2;

create table if not exists chat_room
(
    room     int auto_increment primary key,
    name     varchar(100) not null,
    semester varchar(3)   null,
    constraint `Index 2` unique (semester, name)
) collate = utf8mb4_unicode_ci
  auto_increment = 1724450;

create table if not exists crontab
(
    cron       int auto_increment,
    `interval` int default 7200                                                                                                not null,
    lastRun    int default 0                                                                                                   not null,
    type       enum ('news', 'mensa', 'chat', 'kino', 'roomfinder', 'ticketsale', 'alarm', 'fileDownload', 'canteenHeadCount') null,
    id         int                                                                                                             null,
    constraint cron unique (cron)
) collate = utf8mb4_unicode_ci
  auto_increment = 44;

create table if not exists curricula
(
    curriculum int auto_increment primary key,
    category   enum ('bachelor', 'master') default 'bachelor' not null,
    name       mediumtext                                     not null,
    url        mediumtext                                     not null
) collate = utf8mb4_unicode_ci
  auto_increment = 16;

create table if not exists dish
(
    dish int auto_increment primary key,
    name varchar(150) not null,
    type varchar(20)  not null
) collate = utf8mb4_unicode_ci;

create table if not exists dishflags
(
    flag        int auto_increment primary key,
    short       varchar(10) not null,
    description varchar(50) not null
) collate = utf8mb4_unicode_ci;

create table if not exists dish2dishflags
(
    dish2dishflags int auto_increment primary key,
    dish           int not null,
    flag           int not null,
    constraint dish unique (dish, flag),
    constraint dish2dishflags_ibfk_1 foreign key (dish) references dish (dish) on update cascade on delete cascade,
    constraint dish2dishflags_ibfk_2 foreign key (flag) references dishflags (flag) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci;

create or replace index flag on dish2dishflags (flag);

create table if not exists faculty
(
    faculty int auto_increment primary key,
    name    varchar(150) not null,
    constraint name unique (name)
) charset = utf8mb4
  auto_increment = 18;

create table if not exists feedback
(
    id          int auto_increment primary key,
    email_id    text charset utf8                                      null,
    receiver    text charset utf8                                      null,
    reply_to    text charset utf8                                      null,
    feedback    text charset utf8                                      null,
    image_count int                                                    null,
    latitude    decimal(11, 8)                                         null,
    longitude   decimal(11, 8)                                         null,
    timestamp   datetime /* mariadb-5.3 */ default current_timestamp() null
) auto_increment = 293;

create table if not exists files
(
    file       int auto_increment primary key,
    name       mediumtext           not null,
    path       mediumtext           not null,
    downloads  int        default 0 not null,
    url        varchar(191)         null,
    downloaded tinyint(1) default 1 null
) collate = utf8mb4_unicode_ci
  auto_increment = 34761;

create table if not exists kino
(
    kino        int auto_increment primary key,
    date        datetime /* mariadb-5.3 */                              not null,
    created     timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    title       text                                                    not null,
    year        varchar(4)                                              not null,
    runtime     varchar(40)                                             not null,
    genre       varchar(100)                                            not null,
    director    text                                                    not null,
    actors      text                                                    not null,
    rating      varchar(4)                                              not null,
    description text                                                    not null,
    cover       int                                                     null,
    trailer     text                                                    null,
    link        varchar(190)                                            not null,
    constraint link unique (link),
    constraint kino_ibfk_1 foreign key (cover) references files (file) on update cascade on delete set null
) collate = utf8mb4_unicode_ci
  auto_increment = 219;

create or replace index cover on kino (cover);

create table if not exists lecture
(
    lecture int auto_increment primary key,
    title   varchar(255) null
);

create table if not exists location
(
    location int auto_increment primary key,
    name     text             not null,
    lon      float(10, 6)     not null,
    lat      float(10, 6)     not null,
    radius   int default 1000 not null comment 'in meters'
) charset = utf8
  auto_increment = 2;

create table if not exists member
(
    member          int auto_increment primary key,
    lrz_id          varchar(7)    not null,
    name            varchar(150)  not null,
    active_day      int default 0 null,
    active_day_date date          null,
    student_id      text          null,
    employee_id     text          null,
    external_id     text          null,
    constraint lrz_id unique (lrz_id)
) collate = utf8mb4_unicode_ci
  auto_increment = 104353;

create table if not exists card
(
    card           int auto_increment primary key,
    member         int                  null,
    lecture        int                  null,
    card_type      int                  null,
    title          varchar(255)         null,
    front_text     varchar(2000)        null,
    front_image    varchar(2000)        null,
    back_text      varchar(2000)        null,
    back_image     varchar(2000)        null,
    can_shift      tinyint(1) default 0 null,
    created_at     date                 not null,
    updated_at     date                 not null,
    duplicate_card int                  null,
    aggr_rating    float      default 0 null,
    constraint card_ibfk_1 foreign key (member) references member (member) on delete set null,
    constraint card_ibfk_2 foreign key (lecture) references lecture (lecture) on delete set null,
    constraint card_ibfk_3 foreign key (card_type) references card_type (card_type) on delete set null,
    constraint card_ibfk_4 foreign key (duplicate_card) references card (card) on delete set null
);

create or replace index card_type on card (card_type);

create or replace index duplicate_card on card (duplicate_card);

create or replace index lecture on card (lecture);

create or replace index member on card (member);

create table if not exists card_box
(
    card_box int auto_increment primary key,
    member   int          null,
    title    varchar(255) null,
    duration int          not null,
    constraint card_box_ibfk_1 foreign key (member) references member (member) on delete cascade
) auto_increment = 6;

create or replace index member on card_box (member);

create table if not exists card_comment
(
    card_comment int auto_increment primary key,
    member       int           null,
    card         int           not null,
    rating       int default 0 null,
    created_at   date          not null,
    constraint card_comment_ibfk_1 foreign key (member) references member (member) on delete set null,
    constraint card_comment_ibfk_2 foreign key (card) references card (card) on delete cascade
);

create or replace index card on card_comment (card);

create or replace index member on card_comment (member);

create table if not exists card_option
(
    card_option       int auto_increment primary key,
    card              int                      not null,
    text              varchar(2000) default '' null,
    is_correct_answer tinyint(1)    default 0  null,
    sort_order        int           default 0  not null,
    image             varchar(2000)            null,
    constraint card_option_ibfk_1 foreign key (card) references card (card) on delete cascade
);

create or replace index card on card_option (card);

create table if not exists chat_message
(
    message   int auto_increment primary key,
    member    int                        not null,
    room      int                        not null,
    text      longtext                   not null,
    created   datetime /* mariadb-5.3 */ not null,
    signature longtext                   not null,
    constraint FK_chat_message_chat_room foreign key (room) references chat_room (room) on update cascade on delete cascade,
    constraint chat_message_ibfk_1 foreign key (member) references member (member) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci
  auto_increment = 1977;

create or replace index chat_message_b3c09425 on chat_message (member);

create or replace index chat_message_ca20ebca on chat_message (room);

create table if not exists chat_room2members
(
    room2members int auto_increment primary key,
    room         int not null,
    member       int not null,
    constraint chatroom_id unique (room, member),
    constraint FK_chat_room2members_chat_room foreign key (room) references chat_room (room) on update cascade on delete cascade,
    constraint chat_room2members_ibfk_2 foreign key (member) references member (member) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci
  auto_increment = 63377426;

create or replace index chat_chatroom_members_29801a33 on chat_room2members (room);

create or replace index chat_chatroom_members_b3c09425 on chat_room2members (member);

create table if not exists devices
(
    device          int auto_increment primary key,
    member          int                                                       null,
    uuid            varchar(50)                                               not null,
    created         timestamp /* mariadb-5.3 */                               null,
    lastAccess      timestamp /* mariadb-5.3 */ default '0000-00-00 00:00:00' not null on update current_timestamp(),
    lastApi         mediumtext                  default ''                    not null,
    developer       enum ('true', 'false')      default 'false'               not null,
    osVersion       mediumtext                  default ''                    not null,
    appVersion      mediumtext                  default ''                    not null,
    counter         int                         default 0                     not null,
    pk              longtext                                                  null,
    pkActive        enum ('true', 'false')      default 'false'               not null,
    gcmToken        text                                                      null,
    gcmStatus       varchar(200)                                              null,
    confirmationKey varchar(35)                                               null,
    keyCreated      datetime                                                  null,
    keyConfirmed    datetime                                                  null,
    constraint uuid unique (uuid),
    constraint devices_ibfk_1 foreign key (member) references member (member) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci
  auto_increment = 144352;

create table if not exists device2stats
(
    device                        int default 0 not null primary key,
    LecturesPersonalActivity      int default 0 not null,
    CafeteriaActivity             int default 0 not null,
    WizNavStartActivity           int default 0 not null,
    NewsActivity                  int default 0 not null,
    StartupActivity               int default 0 not null,
    MainActivity                  int default 0 not null,
    ChatRoomsActivity             int default 0 not null,
    CalendarActivity              int default 0 not null,
    WizNavCheckTokenActivity      int default 0 not null,
    ChatActivity                  int default 0 not null,
    CurriculaActivity             int default 0 not null,
    CurriculaDetailsActivity      int default 0 not null,
    GradeChartActivity            int default 0 not null,
    GradesActivity                int default 0 not null,
    InformationActivity           int default 0 not null,
    LecturesAppointmentsActivity  int default 0 not null,
    LecturesDetailsActivity       int default 0 not null,
    OpeningHoursDetailActivity    int default 0 not null,
    OpeningHoursListActivity      int default 0 not null,
    OrganisationActivity          int default 0 not null,
    OrganisationDetailsActivity   int default 0 not null,
    PersonsDetailsActivity        int default 0 not null,
    PersonsSearchActivity         int default 0 not null,
    PlansActivity                 int default 0 not null,
    PlansDetailsActivity          int default 0 not null,
    RoomFinderActivity            int default 0 not null,
    RoomFinderDetailsActivity     int default 0 not null,
    SetupEduroamActivity          int default 0 not null,
    TransportationActivity        int default 0 not null,
    TransportationDetailsActivity int default 0 not null,
    TuitionFeesActivity           int default 0 not null,
    UserPreferencesActivity       int default 0 not null,
    WizNavExtrasActivity          int default 0 not null,
    WizNavChatActivity            int default 0 not null,
    TuitionFeesCard               int default 0 not null,
    NextLectureCard               int default 0 not null,
    CafeteriaMenuCard             int default 0 not null,
    NewsCard1                     int default 0 not null,
    NewsCard2                     int default 0 not null,
    NewsCard3                     int default 0 not null,
    NewsCard7                     int default 0 not null,
    constraint device2stats_ibfk_2 foreign key (device) references devices (device) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci;

create or replace index member on devices (member);

create table if not exists members_card
(
    member                   int not null,
    card                     int not null,
    card_box                 int null,
    last_answered_active_day int null,
    primary key (member, card),
    constraint members_card_ibfk_1 foreign key (member) references member (member) on delete cascade,
    constraint members_card_ibfk_2 foreign key (card) references card (card) on delete cascade,
    constraint members_card_ibfk_3 foreign key (card_box) references card_box (card_box) on delete set null
);

create or replace index card on members_card (card);

create or replace index card_box on members_card (card_box);

create table if not exists members_card_answer_history
(
    members_card_answer_history int auto_increment primary key,
    member                      int                       not null,
    card                        int                       null,
    card_box                    int                       null,
    answer                      varchar(2000)             null,
    answer_score                float(10, 2) default 0.00 null,
    created_at                  date                      not null,
    constraint members_card_answer_history_ibfk_1 foreign key (member) references member (member) on delete cascade,
    constraint members_card_answer_history_ibfk_2 foreign key (card) references card (card) on delete set null,
    constraint members_card_answer_history_ibfk_3 foreign key (card_box) references card_box (card_box) on delete set null
);

create or replace index card on members_card_answer_history (card);

create or replace index card_box on members_card_answer_history (card_box);

create or replace index member on members_card_answer_history (member);

create table if not exists mensa
(
    mensa     int auto_increment primary key,
    id        int                           null,
    name      mediumtext                    not null,
    address   mediumtext                    not null,
    latitude  float(10, 6) default 0.000000 not null,
    longitude float(10, 6) default 0.000000 not null,
    constraint id unique (id)
) collate = utf8mb4_unicode_ci
  auto_increment = 17;

create table if not exists dish2mensa
(
    dish2mensa int auto_increment primary key,
    mensa      int                                                       not null,
    dish       int                                                       not null,
    date       date                                                      not null,
    created    datetime /* mariadb-5.3 */                                not null,
    modifierd  timestamp /* mariadb-5.3 */ default '0000-00-00 00:00:00' not null on update current_timestamp(),
    constraint dish2mensa_ibfk_1 foreign key (mensa) references mensa (mensa) on update cascade on delete cascade,
    constraint dish2mensa_ibfk_2 foreign key (dish) references dish (dish) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci;

create or replace index dish on dish2mensa (dish);

create or replace index mensa on dish2mensa (mensa);

create table if not exists mensaplan_mensa
(
    id        int auto_increment primary key,
    name      varchar(100) not null,
    latitude  double       null,
    longitude double       null,
    webid     int          null,
    category  varchar(50)  not null
) engine = MyISAM
  collate = utf8mb4_unicode_ci
  auto_increment = 30;

create table if not exists mensaprices
(
    price      int auto_increment primary key,
    created    timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    person     mediumtext                                              not null,
    type       mediumtext                                              not null,
    typeLong   mediumtext                                              not null,
    typeNumber int                                                     not null,
    value      decimal                                                 not null
) collate = utf8mb4_unicode_ci;

create table if not exists migrations
(
    id varchar(255) not null primary key
);

create table if not exists newsSource
(
    source int auto_increment primary key,
    title  mediumtext                         not null,
    url    mediumtext                         null,
    icon   int                                not null,
    hook   enum ('newspread', 'impulsivHook') null,
    constraint newsSource_ibfk_1 foreign key (icon) references files (file) on update cascade on delete set null
) collate = utf8mb4_unicode_ci
  auto_increment = 17;

create table if not exists news
(
    news        int auto_increment primary key,
    date        datetime /* mariadb-5.3 */                              not null,
    created     timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    title       tinytext                                                not null,
    description text                                                    not null,
    src         int                                                     not null,
    link        varchar(190)                                            not null,
    image       text                                                    null,
    file        int                                                     null,
    constraint link unique (link),
    constraint news_ibfk_1 foreign key (src) references newsSource (source) on update cascade on delete cascade,
    constraint news_ibfk_2 foreign key (file) references files (file) on update cascade on delete set null
) collate = utf8mb4_unicode_ci
  auto_increment = 770113;

create or replace index file on news (file);

create or replace index src on news (src);

create or replace index icon on newsSource (icon);

create table if not exists news_alert
(
    news_alert int auto_increment primary key,
    file       int                                                     null,
    name       varchar(100)                                            null,
    link       text                                                    null,
    created    timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    `from`     datetime /* mariadb-5.3 */  default current_timestamp() not null,
    `to`       datetime /* mariadb-5.3 */  default current_timestamp() not null,
    constraint FK_File unique (file)
) charset = utf8mb4
  auto_increment = 7;

create table if not exists notification_type
(
    type         int auto_increment primary key,
    name         text                                   not null,
    confirmation enum ('true', 'false') default 'false' not null
) charset = utf8
  auto_increment = 4;

create table if not exists notification
(
    notification int auto_increment primary key,
    type         int                                                     not null,
    location     int                                                     null,
    title        text                                                    not null,
    description  text                                                    not null,
    created      timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    signature    text                                                    null,
    silent       tinyint(1)                  default 0                   not null,
    constraint notification_ibfk_1 foreign key (type) references notification_type (type) on update cascade on delete cascade,
    constraint notification_ibfk_2 foreign key (location) references location (location) on update cascade on delete set null
) charset = utf8
  auto_increment = 107;

create or replace index location on notification (location);

create or replace index type on notification (type);

create table if not exists notification_confirmation
(
    notification int                                                     not null,
    device       int                                                     not null,
    sent         tinyint(1)                  default 0                   not null,
    created      timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    received     timestamp /* mariadb-5.3 */                             null,
    primary key (notification, device),
    constraint notification_confirmation_ibfk_1 foreign key (notification) references notification (notification),
    constraint notification_confirmation_ibfk_2 foreign key (device) references devices (device)
) charset = utf8;

create or replace index device on notification_confirmation (device);

create table if not exists openinghours
(
    id                int auto_increment primary key,
    category          varchar(20)             not null,
    name              varchar(60)             not null,
    address           varchar(140)            not null,
    room              varchar(140)            null,
    transport_station varchar(150)            null,
    opening_hours     varchar(300)            null,
    infos             varchar(500)            null,
    url               varchar(300)            not null,
    language          varchar(2) default 'de' null,
    reference_id      int        default -1   null
) collate = utf8mb4_unicode_ci
  auto_increment = 113;

create table if not exists question
(
    question int auto_increment primary key,
    member   int                                                     not null,
    text     text                                                    not null,
    created  timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    end      timestamp /* mariadb-5.3 */                             null
) auto_increment = 282;

create or replace index member on question (member);

create table if not exists question2answer
(
    question int not null,
    answer   int not null,
    member   int not null,
    constraint question unique (question, member)
);

create table if not exists question2faculty
(
    question int not null,
    faculty  int not null,
    primary key (question, faculty),
    constraint question2faculty_ibfk_1 foreign key (question) references question (question) on update cascade on delete cascade,
    constraint question2faculty_ibfk_2 foreign key (faculty) references faculty (faculty) on update cascade on delete cascade
);

create or replace index faculty on question2faculty (faculty);

create table if not exists questionAnswers
(
    answer int auto_increment primary key,
    text   text not null
) auto_increment = 3;

create table if not exists reports
(
    report             int auto_increment primary key,
    device             int                                                     null,
    created            timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    fixed              enum ('true', 'false')      default 'false'             not null,
    issue              int                                                     null,
    stacktrace         mediumtext                                              not null,
    stacktraceGroup    text                                                    null,
    log                mediumtext                                              not null,
    package            mediumtext                                              not null,
    packageVersion     mediumtext                                              not null,
    packageVersionCode int                         default -1                  not null,
    model              mediumtext                                              not null,
    osVersion          mediumtext                                              not null,
    networkWifi        varchar(10)                                             not null,
    networkMobile      varchar(10)                                             not null,
    gps                varchar(10)                                             not null,
    screenWidth        varchar(100)                                            not null,
    screenHeight       varchar(100)                                            not null,
    screenOrientation  varchar(100)                                            not null,
    screenDpi          varchar(100)                                            not null,
    constraint reports_ibfk_3 foreign key (device) references devices (device) on update cascade on delete set null
) collate = utf8mb4_unicode_ci;

create or replace index device on reports (device);

create table if not exists rights
(
    `right`     int auto_increment primary key,
    name        varchar(100)  null,
    description mediumtext    not null,
    category    int default 0 not null,
    constraint Unquie unique (name)
) collate = utf8mb4_unicode_ci
  auto_increment = 14;

create table if not exists menu
(
    menu     int auto_increment primary key,
    `right`  int                                                         null,
    parent   int                                                         null,
    name     varchar(255)                                                null,
    path     varchar(255)                                                null,
    target   enum ('_blank', '_self', '_parent', '_top') default '_self' not null,
    icon     varchar(200)                                                not null,
    class    varchar(200)                                                not null,
    position int                                         default 0       not null,
    constraint menu_ibfk_1 foreign key (`right`) references rights (`right`) on update cascade on delete set null,
    constraint menu_ibfk_2 foreign key (parent) references menu (menu) on update cascade on delete set null
) collate = utf8mb4_unicode_ci
  auto_increment = 25;

create or replace index parent on menu (parent);

create or replace index `right` on menu (`right`);

create table if not exists modules
(
    module  int auto_increment primary key,
    name    varchar(255) null,
    `right` int          null,
    constraint fkMod2Rights foreign key (`right`) references rights (`right`) on update cascade on delete set null
) collate = utf8mb4_unicode_ci
  auto_increment = 31;

create or replace index module_right on modules (`right`);

create table if not exists roles
(
    role        int auto_increment primary key,
    name        varchar(50) not null,
    description mediumtext  not null
) collate = utf8mb4_unicode_ci
  auto_increment = 6;

create table if not exists roles2rights
(
    role    int not null,
    `right` int not null,
    primary key (role, `right`),
    constraint fkRight foreign key (`right`) references rights (`right`) on delete cascade,
    constraint fkRole foreign key (role) references roles (role) on delete cascade
) collate = utf8mb4_unicode_ci;

create or replace index fkRight_idx on roles2rights (`right`);

create table if not exists roomfinder_building2area
(
    area_id     int         not null,
    building_nr varchar(8)  not null primary key,
    campus      char        not null,
    name        varchar(32) not null
) charset = utf8mb4;

grant select on table roomfinder_building2area to 'tca-students';

create table if not exists roomfinder_buildings
(
    building_nr    varchar(8)  not null primary key,
    utm_zone       varchar(4)  null,
    utm_easting    varchar(32) null,
    utm_northing   varchar(32) null,
    default_map_id int         null
) charset = utf8mb4;

grant select on table roomfinder_buildings to 'tca-students';

create table if not exists roomfinder_buildings2gps
(
    id        varchar(8) default '' not null primary key,
    latitude  varchar(30)           null,
    longitude varchar(30)           null
) charset = utf8mb4;

create table if not exists roomfinder_buildings2maps
(
    building_nr varchar(8) not null,
    map_id      int        not null,
    primary key (building_nr, map_id)
) charset = utf8mb4;

grant select on table roomfinder_buildings2maps to 'tca-students';

create table if not exists roomfinder_maps
(
    map_id      int         not null primary key,
    description varchar(64) not null,
    scale       int         not null,
    width       int         not null,
    height      int         not null
) charset = utf8mb4;

grant select on table roomfinder_maps to 'tca-students';

create table if not exists roomfinder_rooms
(
    room_id        int          not null primary key,
    room_code      varchar(32)  null,
    building_nr    varchar(8)   null,
    arch_id        varchar(16)  null,
    info           varchar(64)  null,
    address        varchar(128) null,
    purpose_id     int          null,
    purpose        varchar(64)  null,
    seats          int          null,
    utm_zone       varchar(4)   null,
    utm_easting    varchar(32)  null,
    utm_northing   varchar(32)  null,
    unit_id        int          null,
    default_map_id int          null
) charset = utf8mb4;

grant select on table roomfinder_rooms to 'tca-students';

create table if not exists roomfinder_rooms2maps
(
    room_id int not null,
    map_id  int not null,
    primary key (room_id, map_id)
) charset = utf8mb4;

create table if not exists roomfinder_schedules
(
    room_id     int                        not null,
    start       datetime /* mariadb-5.3 */ not null,
    end         datetime /* mariadb-5.3 */ not null,
    title       varchar(64)                not null,
    event_id    int                        not null,
    course_code varchar(32)                null,
    constraint `unique` unique (room_id, start, end)
) charset = utf8mb4;

grant select on table roomfinder_schedules to 'tca-students';

create table if not exists sessions
(
    session varchar(255) charset utf8 not null primary key,
    access  int unsigned              null,
    data    text                      null,
    constraint session unique (session)
) collate = utf8mb4_unicode_ci;

create table if not exists tag
(
    tag   int auto_increment primary key,
    title varchar(255) null
);

create table if not exists card2tag
(
    tag  int not null,
    card int not null,
    primary key (tag, card),
    constraint card2tag_ibfk_1 foreign key (tag) references tag (tag) on delete cascade,
    constraint card2tag_ibfk_2 foreign key (card) references card (card) on delete cascade
);

create or replace index card on card2tag (card);

create table if not exists ticket_admin
(
    ticket_admin int auto_increment primary key,
    `key`        text                                                    not null,
    created      timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    active       tinyint(1)                  default 0                   not null,
    comment      text                                                    null
) charset = utf8
  auto_increment = 48;

create table if not exists ticket_group
(
    ticket_group int auto_increment primary key,
    description  text not null
) charset = utf8
  auto_increment = 2;

create table if not exists event
(
    event        int auto_increment primary key,
    news         int                        null,
    kino         int                        null,
    file         int                        null,
    title        varchar(100)               not null,
    description  text                       not null,
    locality     varchar(200)               not null,
    link         varchar(200)               null,
    start        datetime /* mariadb-5.3 */ null,
    end          datetime /* mariadb-5.3 */ null,
    ticket_group int default 1              null,
    constraint fkEventFile foreign key (file) references files (file) on update cascade on delete set null,
    constraint fkEventGroup foreign key (ticket_group) references ticket_group (ticket_group),
    constraint fkKino foreign key (kino) references kino (kino) on update cascade on delete set null,
    constraint fkNews foreign key (news) references news (news) on update cascade on delete set null
) charset = utf8
  auto_increment = 39;

create or replace index file on event (file);

create or replace fulltext index searchTitle on event (title);

create table if not exists ticket_admin2group
(
    ticket_admin2group int auto_increment primary key,
    ticket_admin       int not null,
    ticket_group       int not null,
    constraint fkTicketAdmin foreign key (ticket_admin) references ticket_admin (ticket_admin) on update cascade on delete cascade,
    constraint fkTicketGroup foreign key (ticket_group) references ticket_group (ticket_group) on update cascade on delete cascade
) charset = utf8
  auto_increment = 10;

create or replace index ticket_admin on ticket_admin2group (ticket_admin);

create or replace index ticket_group on ticket_admin2group (ticket_group);

create table if not exists ticket_payment
(
    ticket_payment int auto_increment primary key,
    name           varchar(50) not null,
    min_amount     int         null,
    max_amount     int         null,
    config         text        not null
) charset = utf8
  auto_increment = 3;

create table if not exists ticket_type
(
    ticket_type    int auto_increment primary key,
    event          int          not null,
    ticket_payment int          not null,
    price          double       not null,
    contingent     int          not null,
    description    varchar(100) not null,
    constraint fkEvent foreign key (event) references event (event) on update cascade on delete cascade,
    constraint fkPayment foreign key (ticket_payment) references ticket_payment (ticket_payment) on update cascade
) charset = utf8
  auto_increment = 57;

create table if not exists ticket_history
(
    ticket_history int auto_increment primary key,
    member         int                                                     not null,
    ticket_payment int                                                     null,
    ticket_type    int                                                     not null,
    purchase       datetime /* mariadb-5.3 */                              null,
    redemption     datetime /* mariadb-5.3 */                              null,
    created        timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    code           char(128)                                               not null,
    constraint fkMember foreign key (member) references member (member) on update cascade on delete cascade,
    constraint fkTicketPayment foreign key (ticket_payment) references ticket_payment (ticket_payment) on update cascade,
    constraint fkTicketType foreign key (ticket_type) references ticket_type (ticket_type) on update cascade
) charset = utf8
  auto_increment = 776;

create or replace index member on ticket_history (member);

create or replace index ticket_payment on ticket_history (ticket_payment);

create or replace index ticket_type on ticket_history (ticket_type);

create or replace index event on ticket_type (event);

create or replace index ticket_payment on ticket_type (ticket_payment);

create table if not exists update_note
(
    version_code int  not null primary key,
    version_name text null,
    message      text null
);

create table if not exists users
(
    user             int auto_increment primary key,
    username         varchar(7)                                              null,
    firstname        varchar(100)                                            null,
    surname          varchar(100)                                            null,
    created          timestamp /* mariadb-5.3 */ default current_timestamp() not null,
    deleted          int                         default 0                   not null,
    lastActive       int                         default 0                   not null,
    lastPage         text                                                    null,
    lastLogin        datetime                                                null,
    tum_id_student   varchar(50)                                             null comment 'OBFUSCATED_ID_ST',
    tum_id_employee  varchar(50)                                             null comment 'OBFUSCATED_ID_B',
    tum_id_alumni    varchar(50)                                             null comment 'OBFUSCATED_ID_EXT',
    tum_id_preferred varchar(50)                                             null comment 'OBFUSCATED_ID_BEVORZUGT',
    constraint username unique (username)
) collate = utf8mb4_unicode_ci
  auto_increment = 434;

create table if not exists log
(
    log           int auto_increment primary key,
    time          int          not null,
    user_executed int          null,
    user_affected int          null,
    action        int          null,
    comment       varchar(255) not null,
    constraint fkLog2Actions foreign key (action) references actions (action) on update set null on delete set null,
    constraint fkLog2UsersAf foreign key (user_affected) references users (user) on update set null on delete set null,
    constraint fkLog2UsersEx foreign key (user_executed) references users (user) on update set null on delete set null
) collate = utf8mb4_unicode_ci;

create or replace index action on log (action);

create or replace index user on log (user_executed);

create or replace index user_affected on log (user_affected);

create table if not exists recover
(
    recover int auto_increment primary key,
    user    int          not null,
    created int          not null,
    hash    varchar(190) not null,
    ip      varchar(255) not null,
    constraint hash unique (hash),
    constraint fkRecover2User foreign key (user) references users (user) on delete cascade
) collate = utf8mb4_unicode_ci;

create or replace index user on recover (user);

create table if not exists users2info
(
    user         int            not null primary key,
    firstname    varchar(255)   not null,
    surname      varchar(255)   not null,
    lastPwChange int            not null,
    pager        int default 15 null,
    constraint fkUsers foreign key (user) references users (user) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci;

create table if not exists users2roles
(
    user int not null,
    role int not null,
    primary key (user, role),
    constraint fkUser2RolesRole foreign key (role) references roles (role) on update cascade on delete cascade,
    constraint fkUser2RolesUser foreign key (user) references users (user) on update cascade on delete cascade
) collate = utf8mb4_unicode_ci;

create table if not exists wifi_measurement
(
    id               int(20) unsigned auto_increment primary key,
    date             date        not null,
    SSID             varchar(32) not null,
    BSSID            varchar(64) not null,
    dBm              int         null,
    accuracyInMeters float       not null,
    latitude         double      not null,
    longitude        double      not null
) collate = utf8mb4_unicode_ci;
