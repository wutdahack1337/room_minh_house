# introductory-capstone
Tranh thủ đợt niên luận cơ sở này, viết để sau này xài luôn =DDD
- [Slides](https://docs.google.com/presentation/d/1ah9UJaAR0LRt9XVeBSRe9aSXfINHvsAoaEpz0z5HbB8)
- [Notes](https://www.notion.so/Ni-n-lu-n-c-s-ng-nh-KHMT-CT201-2e3d7839804a804383d0d76b004c5f2b)
- [Github](https://github.com/wutdahack1337/introductory-capstone)

---
# Hệ thống quản lý nhà trọ
Trong trường hợp phải quản lý số lượng phòng trọ quá lớn, việc quản lý, ghi chép, tính toán truyền thống tỏ ra không còn mấy hiệu quả. Vì lý do đó, hệ thống này được sinh ra nhằm mục đích hỗ trợ, giảm công sức của người quản lý trọ.

**Người có nhu cầu thuê phòng trọ** thì có khả năng kiểm tra, xem phòng, xem giá, đặt cọc. **Người đã thuê trọ** thì sẽ được người quản lý cấp tài khoản, có khả năng báo cáo sự cố (hư bóng đèn,…) và sẽ nhận được hóa đơn cần được thanh toán hàng tháng. Yêu cầu người thuê trọ phải cập nhật thông tin cá nhân một cách đầy đủ và chính xác.

**Người quản lý** sẽ có khả năng tạo, quản lý tài khoản, tạo, sửa bài cho thuê phòng. Xem tổng kết doanh thu tháng, số phòng còn trống. Tính tiền phòng (bao gồm điện, nước, dịch vụ,…) và xuất, gửi hóa đơn cho tài khoản của người đang thuê phòng.

# Yêu cầu
## Chức năng chính
![alt text](use_case.drawio.png)

## Phi chức năng
Bảo vệ (mã hóa,...) dữ liệu của người thuê trọ

# Công nghệ
- Frontend: Vue, TypeScript
- Backend: Go
- Database: MongoDB

---
Một chủ trọ cần quản lý dãy nhà trọ của mình, mỗi phòng chỉ đủ cho 1 người ở. Một người cùng lúc chỉ có thể thuê 1 phòng. Nhà trọ này chỉ tính tiền phòng, điện và nước, các dịch vụ khác đều là miễn phí.
Hệ thống quản lý nhà trọ cho phép chủ trọ có thể quản lý (thêm, sửa, xóa) phòng trọ, post, người thuê, hợp đồng và hóa đơn. Hệ thống sẽ có trang tổng hợp số liệu, ví dụ: số phòng trống, hợp đồng sắp hết hạn,…

Entities chính sẽ là phòng và người thuê.
Mỗi người thuê có mã số người thuê, số chứng minh nhân dân, số điện thoại (username) và mật khẩu tài khoản. Mỗi người ai cũng sẽ có một thẻ căn cước công dân có số chứng minh nhân dân, họ tên, quê quán và ngày tháng năm sinh. Mỗi người thuê sẽ được chủ trọ tạo và cấp mật khẩu tài khoản sau khi kí hợp đồng. Khi chưa được cấp tài khoản, người dùng chỉ có thể truy cập được trang xem phòng, trang này list ra nhiều phòng, người dùng muốn xem chi tiết phòng nào thì click vào. Sau khi kí hợp đồng và đã được cấp tài khoản, tài khoản sẽ được truy cập thêm vào trang hóa đơn chưa thanh toán và trang thông tin cá nhân. Người thuê trọ bắt buộc phải nhập thông tin cá nhân. Hàng tháng, người thuê sẽ được thông báo về hóa đơn của mình.
Mỗi phòng có mã số phòng, giá thuê mỗi tháng, diện tích, chỉ số điện tháng trước, chỉ số nước tháng trước và tình trạng (đã được thuê chưa). Mỗi phòng có một bài post bao gồm mã bài viết, tựa đề, mô tả, hình ảnh và thời gian tạo được đăng lên trang xem phòng. Mỗi phòng chỉ có duy nhất một bài post. Người dùng có thể sử dụng trang xem phòng, để tham khảo hình ảnh căn phòng, giá phòng, diện tích và dịch vụ đi kèm,…
Hợp đồng có mã số hợp đồng, mã số người thuê, mã phòng, ngày tạo hợp đồng và thời gian hiệu lực. Thời gian hiệu lực mặc định là 3 tháng hoặc 6 tháng. Khi 1 hợp đồng hết hạn thì người thuê sẽ phải kí tạo hợp đồng mới nếu muốn thuê tiếp tục. Người thuê có thể vào trang hợp đồng để xem chi tiết hợp đồng.
Hóa đơn có mã số hóa đơn, mã số phòng, mã số người thuê, tiền điện, tiền nước, ngày tạo và trạng thái (đã thanh toán hay chưa?). Cứ mỗi tháng, tùy vào ngày tạo hợp đồng, hóa đơn sẽ được tự động tạo mới cho mỗi phòng. Hóa đơn chỉ còn thiếu phần chỉ số điện tháng này và chỉ số nước tháng sẽ được nhập bởi chủ trọ. Sau khi chủ trọ hoàn thành việc nhập chỉ số điện và nước, hóa đơn sẽ tự động tính tổng tiền và gửi hóa đơn vào tài khoản của người thuê (trang hóa đơn chưa thanh toán).