# Mô tả Hệ thống Quản lý Nhà Trọ

## Tổng quan

Hệ thống quản lý nhà trọ cho phép chủ trọ quản lý tập trung toàn bộ dữ liệu vận hành, bao gồm phòng trọ, người thuê, hợp đồng thuê nhà và hóa đơn theo từng tháng. Mục tiêu chính của hệ thống là giúp việc theo dõi thông tin thuê trọ trở nên rõ ràng, nhất quán và hạn chế sai sót khi tính tiền điện nước định kỳ.

Trong phạm vi hiện tại, nhà trọ chỉ tính tiền cho ba khoản chính là tiền phòng, tiền điện và tiền nước. Mỗi phòng chỉ đủ cho 1 người ở và mỗi người thuê chỉ được thuê 1 phòng tại một thời điểm, từ đó mô hình dữ liệu được thiết kế theo hướng đơn giản, dễ kiểm soát và phù hợp với nghiệp vụ thực tế đang triển khai.

## Người dùng

Mỗi người dùng trong hệ thống được xác định duy nhất bằng số chứng minh nhân dân (pid - personal ID). Đây là định danh chính để liên kết người dùng với các dữ liệu khác như hợp đồng và thông tin cá nhân. Người dùng có các thông tin cơ bản gồm pid, số điện thoại và role.

Hệ thống hỗ trợ hai role là admin (chủ trọ) và tenant (người thuê). Bên cạnh đó, mỗi người dùng đều có thông tin thẻ căn cước công dân dưới dạng dữ liệu cá nhân đi kèm, bao gồm họ tên, giới tính, ngày tháng năm sinh và quê quán. Việc lưu trữ đầy đủ các trường này giúp quá trình tạo hợp đồng và tra cứu thông tin người thuê được thuận tiện hơn.

## Phòng trọ

Mỗi phòng trọ được định danh bằng mã số phòng và chứa các thông tin cơ bản gồm giá thuê hàng tháng và tình trạng hiện tại (sẵn sàng cho thuê hoặc đã được thuê). Với cấu trúc này, chủ trọ có thể theo dõi nhanh tình trạng khai thác của từng phòng và dễ dàng cập nhật khi có biến động.

Hệ thống cho phép chủ trọ thực hiện đầy đủ các thao tác thêm, sửa và xóa phòng. Khi dữ liệu phòng thay đổi, các chức năng liên quan đến hợp đồng và hóa đơn vẫn dựa trên mã phòng để đảm bảo tính nhất quán giữa các nhóm dữ liệu.

Để tính toán hóa đơn chính xác hàng tháng, hệ thống lưu trữ chỉ số công tơ điện và nước của từng phòng theo từng kỳ (tháng). Mỗi bản ghi kỳ được xác định bằng mã phòng và kỳ hạn theo định dạng MM-YYYY, đồng thời ghi nhận hai giá trị là chỉ số điện và chỉ số nước tại thời điểm chốt kỳ.

Từ dữ liệu này, hệ thống sẽ tự động tính phần tiêu thụ điện nước của kỳ hiện tại bằng cách lấy hiệu số giữa kỳ hiện tại và kỳ trước. Cách tính này giúp giảm nhập liệu thủ công và hạn chế sai lệch trong quá trình lập hóa đơn.

## Hợp đồng thuê nhà

Hợp đồng là dữ liệu ràng buộc giữa một người thuê và một phòng cụ thể. Mỗi hợp đồng được xác định bằng mã hợp đồng duy nhất, đồng thời liên kết trực tiếp với pid của người thuê và mã phòng được thuê. Nhờ đó, hệ thống có thể truy vết rõ ràng mối quan hệ người thuê - phòng - kỳ hóa đơn.

Thông tin hợp đồng bao gồm ngày bắt đầu, ngày kết thúc (được tính tự động dựa trên kỳ hạn) và trạng thái hợp đồng (hoạt động, hết hạn, bị chấm dứt). Việc tính tự động ngày kết thúc giúp giảm thao tác thủ công và đảm bảo logic thời hạn luôn thống nhất.

Khi một hợp đồng hết hạn, nếu người thuê muốn tiếp tục ở, họ phải kí hợp đồng mới.

## Công tơ điện nước

Hàng tháng, chủ trọ cần nhập chỉ số điện và chỉ số nước của từng phòng theo đúng kỳ hạn MM-YYYY. Đây là dữ liệu đầu vào quan trọng để hệ thống tính toán chi phí điện nước cho hóa đơn tháng tương ứng.

Sau khi có đủ dữ liệu của kỳ hiện tại và kỳ trước, hệ thống tự động tính lượng tiêu thụ thực tế bằng hiệu số giữa hai kỳ. Kết quả này sau đó được dùng trực tiếp trong các dòng chi tiết hóa đơn.

## Hóa đơn

Hóa đơn được tạo theo từng tháng cho mỗi phòng đang có hợp đồng ở trạng thái hoạt động. Mỗi hóa đơn được xác định bằng cặp mã hợp đồng và kỳ hạn (composite key), thay vì dùng một mã hóa đơn độc lập. Cách định danh này giúp dữ liệu hóa đơn bám sát ngữ cảnh hợp đồng và tránh trùng lặp theo kỳ.

Hóa đơn chứa tổng số tiền cần thanh toán và trạng thái thanh toán (chưa thanh toán hoặc đã thanh toán). Trong đó, tổng tiền được tổng hợp từ các dòng chi tiết theo từng loại dịch vụ.

Chi tiết hóa đơn được chia thành ba dòng chuẩn: tiền phòng (cố định theo giá phòng), tiền điện (số khối điện nhân với giá điện), và tiền nước (số khối nước nhân với giá nước). Mỗi dòng chi tiết lưu loại dịch vụ, số lượng tiêu thụ, đơn giá và thành tiền để đảm bảo dễ kiểm tra và đối soát.

Quy trình tạo hóa đơn:
1. Hệ thống tự động tạo khung hóa đơn cho kỳ hạn kế tiếp khi hợp đồng còn hoạt động.
2. Chủ trọ nhập chỉ số điện và nước của tháng đó.
3. Hệ thống tự động tính toán lượng tiêu thụ (hiệu số của hai tháng liên tiếp) và chi tiết từng dòng.

Quy trình trên giúp tách rõ vai trò: chủ trọ chịu trách nhiệm nhập số liệu công tơ, còn hệ thống chịu trách nhiệm tính toán và chuẩn hóa hóa đơn. Nhờ đó, dữ liệu phát sinh theo tháng có thể được theo dõi rõ ràng hơn và giảm rủi ro sai công thức khi tính tiền.

## Giao diện người dùng

**Trang dành cho Admin (chủ trọ):**
- Quản lý phòng: Thêm, sửa, xóa phòng trọ.
- Quản lý người thuê: Thêm, sửa, xóa tài khoản người thuê.
- Quản lý hợp đồng: Tạo, cập nhật, hủy hợp đồng.
- Quản lý hóa đơn: Tạo hóa đơn hàng tháng, nhập chỉ số điện nước, cập nhật trạng thái thanh toán.

Trong trạng thái triển khai hiện tại, giao diện tập trung vào nhóm chức năng quản trị cho admin để phục vụ vận hành. Các màn hình và luồng thao tác đều ưu tiên cho việc nhập liệu, kiểm tra và cập nhật dữ liệu nghiệp vụ cốt lõi, bao gồm phòng, người thuê, hợp đồng và hóa đơn.

## Công nghệ

- **Frontend**: Vue với TypeScript
- **Backend**: Go (gRPC + Gateway REST)
- **Database**: MongoDB
