package response

import (
	"context"
	"fmt"
	pb "golang-bot/pb"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	RPCRequest struct {
		AddressURI string
	}
	RPCResponse struct {
		Connection *grpc.ClientConn
		Close      func()
	}
)

func callRPC(request RPCRequest) (RPCResponse, error) {
	conn, err := grpc.Dial(request.AddressURI, grpc.WithInsecure())
	if err != nil {
		return RPCResponse{}, err
	}

	return RPCResponse{
		Connection: conn,
		Close: func() {
			err = conn.Close()
		},
	}, nil
}

func Help() string {
	strHelp := []string{
		"**!list** - Melihat daftar buku",
		"**!tambah title;year;pages;content** - Menambah buku",
		"**!update id;title;year;pages;content** - Mengupdate buku",
		"**!hapus id** - Menghapus buku berdasarkan id",
	}
	return strings.Join(strHelp, "\n")
}

func GetAllBook() string {
	grpcResponse, err := callRPC(RPCRequest{
		AddressURI: viper.GetString("grpc-server"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "something went wrong, plz try again"
	}
	defer grpcResponse.Close()

	c := pb.NewBookGrpcClient(grpcResponse.Connection)
	res, err := c.GetAllBook(context.Background(), &pb.Message{})
	if err != nil {
		fmt.Println(err.Error())
		return "cannot fetch, try again"
	}

	books := res.GetAllBook()
	var txtStr []string

	txtStr = append(txtStr, fmt.Sprintf("**List Buku** : %d", len(books)))
	for _, v := range books {
		txtStr = append(txtStr, fmt.Sprintf(":white_small_square:  %s\n Judul : %s | Halaman : %d | Tahun : %d\n%s", v.ID, v.Title, v.Pages, v.Year, v.Content))
	}

	return strings.Join(txtStr, "\n")
}

func CreateBook(contentBook string) string {
	grpcResponse, err := callRPC(RPCRequest{
		AddressURI: viper.GetString("grpc-server"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "something went wrong, plz try again"
	}
	defer grpcResponse.Close()

	//title;year;pages;content
	book := strings.Split(contentBook, ";")
	pages, err := strconv.Atoi(book[2])
	if err != nil {
		return "wrong format\npls fill : title;year;pages;content"
	}
	year, err := strconv.Atoi(book[1])
	if err != nil {
		return "wrong format\npls fill : title;year;pages;content"
	}
	c := pb.NewBookGrpcClient(grpcResponse.Connection)
	res, err := c.CreateNewBook(context.Background(), &pb.BookPayload{
		Pages:   int32(pages),
		Year:    int32(year),
		Title:   book[0],
		Content: book[3],
	})
	if err != nil {
		fmt.Println(err.Error())
		return "cannot save book info, try again"
	}

	return res.Message
}

func UpdateBook(contentBook string) string {
	grpcResponse, err := callRPC(RPCRequest{
		AddressURI: viper.GetString("grpc-server"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "something went wrong, plz try again"
	}
	defer grpcResponse.Close()

	//id;title;year;pages;content
	book := strings.Split(contentBook, ";")
	pages, err := strconv.Atoi(book[3])
	if err != nil {
		return "wrong format\npls fill : title;year;pages;content"
	}
	year, err := strconv.Atoi(book[2])
	if err != nil {
		return "wrong format\npls fill : title;year;pages;content"
	}
	c := pb.NewBookGrpcClient(grpcResponse.Connection)
	res, err := c.UpdateByIdBook(context.Background(), &pb.UpdateBook{
		Id: &pb.BookId{
			ID: book[0],
		},
		Book: &pb.BookPayload{
			Pages:   int32(pages),
			Year:    int32(year),
			Title:   book[1],
			Content: book[4],
		},
	})
	if err != nil {
		s := status.Convert(err)
		switch s.Code() {
		case codes.Unknown:
			return "id not found"
		}
		fmt.Println(err.Error())
		return "cannot update book info, try again"
	}

	return res.Message
}

func DeleteBook(contentBook string) string {
	grpcResponse, err := callRPC(RPCRequest{
		AddressURI: viper.GetString("grpc-server"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return "something went wrong, plz try again"
	}
	defer grpcResponse.Close()

	c := pb.NewBookGrpcClient(grpcResponse.Connection)
	res, err := c.DelByIdBook(context.Background(), &pb.BookId{
		ID: contentBook,
	})
	if err != nil {
		s := status.Convert(err)
		switch s.Code() {
		case codes.Unknown:
			return "id not found"
		}
		fmt.Println(err.Error())
		return "cannot delete, try again"
	}

	return res.Message
}
