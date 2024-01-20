package test

import (
	General "chs_cloud_general/internal/general"
	GlobalVar "chs_cloud_general/internal/global_var"
	"strconv"

	"errors"
	"fmt"
	"testing"

	"gorm.io/gorm"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
// func TestHelloName(t *testing.T) {
// 	name := "Gladys"
// 	want := regexp.MustCompile(`\b` + name + `\b`)
// 	msg, err := Hello("Gladys")
// 	if !want.MatchString(msg) || err != nil {
// 		t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
// 	}
// }

// // TestHelloEmpty calls greetings.Hello with an empty string,
// // checking for an error.
func TestHelloEmpty(t *testing.T) {

	roomWidth := 100
	roomHeight := 100
	row := 0
	col := 0
	for i := 0; i < 12; i++ { // Increase the loop count to create a 3x4 grid of positions
		if i%3 == 0 {
			row += 1
			col = 0
		}

		posX := roomWidth * col
		posY := roomHeight * row

		col++

		fmt.Println(posX)
		fmt.Println(posY)
	}
}

func StructToJsonObject(StructData interface{}) {
	// jsonBytes, err := json.Marshal(StructData)
	// if err != nil {
	// 	//fmt.Println(err)
	// 	return
	// }
	// jsonString := string(jsonBytes)
	//fmt.Println(jsonString)
}

func TestConv(t *testing.T) {
	number, err := strconv.Atoi("5.000")
	if err != nil {
		fmt.Println(err.Error())
		number = 0
	}
	fmt.Println(number)
}

func TestFunction(t *testing.T) {
	// AuditDate := GlobalQuery. GetAuditDate(c, DB, false)
	AuditDateStr := "" //General.FormatDate1(AuditDate)
	var DataOutput []map[string]interface{}
	err := GlobalVar.Db.Raw(
		"SELECT" +
			" cfg_init_room.number AS RoomNumber," +
			" cfg_init_room.room_type_code," +
			" cfg_init_room.building," +
			" cfg_init_room.floor," +
			" cfg_init_room.status_code," +
			" cfg_init_room.pos_x," +
			" cfg_init_room.pos_y," +
			" cfg_init_room.width," +
			" cfg_init_room.height," +
			" A.*," +
			" IFNULL(A.is_incognito, '0') AS IsIncognito," +
			" IF(IFNULL(LastCheckOut.guest_profile_id1, 0) = 0, '0', '-1') AS IsRepeater," +
			" CONCAT(cfg_init_room_type.name, ' ', cfg_init_bed_type.name) AS RoomType, " +
			"IF(TypeX=UnavailableMarkX,IF(IFNULL(Number, '')<>'', RoomStatusCode, ''),IF(TypeX=FolioMarkX,  CONCAT(IF(IFNULL(Number, '') <> '', 'O', 'V'), cfg_init_room.status_code), IF(TypeX=FolioMarkX, ISI, IF(TypeX=FolioMarkX, ISI, '')))) AS RoomStatusIconCode " +

			"FROM" +
			" cfg_init_room" +
			" LEFT OUTER JOIN (( " +
			"SELECT" +
			" 'FF' AS TypeX," +
			" folio.number AS Number," +
			" folio.guest_profile_id1," +
			" folio.room_status_code AS OccupiedStatusCode," +
			" folio.compliment_hu AS RoomStatusCode," +
			" folio.is_incognito," +
			" guest_group.name AS GroupName," +
			" CONCAT(contact_person.title_code, contact_person.full_name) AS FullName," +
			" DATE(guest_detail.arrival) AS arrival," +
			" DATE(guest_detail.departure) AS departure," +
			" guest_detail.adult," +
			" guest_detail.child," +
			" guest_detail.room_number," +
			" SUM(IF(sub_folio.type_code='D', IFNULL(sub_folio.quantity * sub_folio.amount, 0), -IFNULL(sub_folio.quantity * sub_folio.amount, 0))) AS Balance," +
			" IFNULL(guest_message.id,0) AS Message," +
			" IFNULL(guest_to_do.id,0) AS ToDo " +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN guest_group ON (folio.group_code = guest_group.code)" +
			" LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN sub_folio ON (folio.number = sub_folio.folio_number AND sub_folio.void='0')" +
			" LEFT OUTER JOIN guest_message ON (folio.number = guest_message.folio_number AND guest_message.is_delivered='0')" +
			" LEFT OUTER JOIN guest_to_do ON (folio.number = guest_to_do.folio_number AND guest_to_do.is_done='0')" +
			" WHERE folio.status_code='" + GlobalVar.FolioStatus.Open + "'" +
			"GROUP BY folio.number" +
			")UNION( " +
			"SELECT" +
			" 'RR' AS TypeX," +
			" reservation.number AS Number," +
			" reservation.guest_profile_id1," +
			" '' AS OccupiedStatusCode," +
			" '' AS RoomStatusCode," +
			" '0' AS is_incognito," +
			" guest_group.name AS GroupName," +
			" CONCAT(contact_person.title_code, contact_person.full_name) AS FullName," +
			" DATE(guest_detail.arrival) AS arrival," +
			" DATE(guest_detail.departure) AS departure," +
			" guest_detail.adult," +
			" guest_detail.child," +
			" guest_detail.room_number," +
			" SUM(CASE WHEN guest_deposit.type_code='D' THEN IFNULL(guest_deposit.amount, 0) else -IFNULL(guest_deposit.amount, 0) END) AS Balance," +
			" 0 AS Message," +
			" 0 AS ToDo " +
			"FROM" +
			" reservation" +
			" LEFT OUTER JOIN guest_group ON (reservation.group_code = guest_group.code)" +
			" LEFT OUTER JOIN contact_person ON (reservation.contact_person_id1 = contact_person.id)" +
			" LEFT OUTER JOIN guest_detail ON (reservation.guest_detail_id = guest_detail.id)" +
			" LEFT OUTER JOIN guest_deposit ON (reservation.number = guest_deposit.reservation_number AND guest_deposit.void='0' AND guest_deposit.system_code='" + GlobalVar.ConstProgramVariable.DefaultSystemCode + "')" +
			" WHERE reservation.status_code='" + GlobalVar.ReservationStatus.New + "'" +
			" AND guest_detail.room_number<>''" +
			" AND DATE(guest_detail.arrival)<='" + AuditDateStr + "'" +
			" AND DATE(guest_detail.departure)>'" + AuditDateStr + "' " +
			"GROUP BY reservation.number" +
			")UNION( " +
			"SELECT" +
			" 'UU' AS TypeX," +
			" room_unavailable.id AS Number," +
			" 0 AS guest_profile_id," +
			" '' AS OccupiedStatusCode," +
			" room_unavailable.status_code AS RoomStatusCode," +
			" '0' AS is_incognito," +
			" '' AS GroupName," +
			" cfg_init_room_unavailable_reason.description AS FullName," +
			" room_unavailable.start_date AS arrival," +
			" room_unavailable.end_date AS departure," +
			" COUNT(0) AS adult," +
			" COUNT(0) AS child," +
			" room_unavailable.room_number," +
			" 0 AS Balance," +
			" 0 AS Message," +
			" 0 AS ToDo " +
			"FROM" +
			" room_unavailable" +
			" LEFT OUTER JOIN cfg_init_room_unavailable_reason ON (room_unavailable.reason_code = cfg_init_room_unavailable_reason.code)" +
			" WHERE room_unavailable.start_date<='" + AuditDateStr + "' AND room_unavailable.end_date>='" + AuditDateStr + "')) AS A ON(cfg_init_room.number = A.room_number)" +
			" LEFT OUTER JOIN (" +
			"SELECT DISTINCT" +
			" folio.guest_profile_id1 " +
			"FROM" +
			" folio" +
			" LEFT OUTER JOIN contact_person ON (folio.contact_person_id1 = contact_person.id)" +
			" LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)" +
			" WHERE folio.status_code='" + GlobalVar.FolioStatus.Closed + "'" +
			" AND DATE(guest_detail.departure)<='" + AuditDateStr + "' " +
			") AS LastCheckOut ON(A.guest_profile_id1 = LastCheckOut.guest_profile_id1)" +
			" LEFT OUTER JOIN cfg_init_room_type ON (cfg_init_room.room_type_code = cfg_init_room_type.code)" +
			" LEFT OUTER JOIN cfg_init_bed_type ON (cfg_init_room.bed_type_code = cfg_init_bed_type.code) " +
			"GROUP BY cfg_init_room.number " +
			"ORDER BY cfg_init_room.building, cfg_init_room.floor, cfg_init_room.number").Scan(&DataOutput).Error
	if err != nil {

		fmt.Println(err.Error())
	}

	fmt.Println(DataOutput)

	// DB.Table(DBVar.TableName.Folio).Select(
	// 	" COUNT(IF(DATE(guest_detail.departure)='"+PostingDateStrTomorrow+"', folio.number, NULL)) AS TotalCOTomorrow,"+
	// 		" SUM(IF(DATE(guest_detail.departure)='"+PostingDateStrTomorrow+"', guest_detail.adult + guest_detail.child, NULL)) AS TotalPersonCOTomorrow ").
	// 	Joins("LEFT OUTER JOIN guest_detail ON (folio.guest_detail_id = guest_detail.id)").
	// 	Where("guest_detail.departure_unixx=UNIX_TIMESTAMP(?)", PostingDateStrTomorrow).
	// 	Where("folio.type_code=?", GlobalVar.FolioType.GuestFolio).
	// 	Where("folio.status_code<>?", GlobalVar.FolioStatus.CancelCheckIn).
	// 	Scan(&Data)

	// fmt.Println(Data)
}

func TestEncrypt(t *testing.T) {

	Encrypt, _ := General.EncryptString(GlobalVar.EncryptKey, "cakratendados")
	fmt.Println(Encrypt)
}

func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}
	// Create a message using a random format.
	// message := fmt.Sprintf(randomFormat(), name)
	message := fmt.Sprint("jik")
	return message, nil
}

type Thing1 struct {
	gorm.Model
	Name string `gorm:"size:20"`
	One  int
}

type Thing2 struct {
	gorm.Model
	Name string `gorm:"size:20"`
	Two  int
}

type Thing3 struct {
	gorm.Model
	Name  string `gorm:"size:20"`
	Three int
}

type Composite struct {
	Thing1
	Thing2 Thing2 `gorm:"foreignKey:ID;references:ID"`
	Thing3 Thing3 `gorm:"foreignKey:ID;references:ID"`
}

func TestGORM(t *testing.T) {
	result := &Composite{}

	if err := GlobalVar.Db.Migrator().DropTable(&Thing1{}, &Thing2{}, &Thing3{}); err != nil {
		t.Errorf("Failed, got error: %v", err)
	}

	if err := GlobalVar.Db.AutoMigrate(&Thing1{}, &Thing2{}, &Thing3{}); err != nil {
		t.Errorf("Failed, got error: %v", err)
	}

	GlobalVar.Db.Create(&Thing1{
		Name: "Thing 1",
		One:  1,
	})
	GlobalVar.Db.Create(&Thing2{
		Name: "Thing 2",
		Two:  2,
	})
	GlobalVar.Db.Create(&Thing3{
		Name:  "Thing 3",
		Three: 3,
	})

	GlobalVar.Db.Table("thing1").Joins("Thing2").Joins("Thing3").Find(result)
	//fmt.Println(result)
}
