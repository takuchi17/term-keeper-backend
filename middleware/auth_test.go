package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/takuchi17/term-keeper/configs"
)

func TestAuthMiddleware(t *testing.T) {
	// テスト用のシークレットキーを設定
	configs.Config.JWTSecret = "test-secret-key"
	jwtSecret = []byte(configs.Config.JWTSecret)

	// テーブル駆動テスト
	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedUserID string
		expectedName   string
	}{
		{
			name:           "有効なトークン",
			token:          generateValidToken(t, "user123", "テストユーザー"),
			expectedStatus: http.StatusOK,
			expectedUserID: "user123",
			expectedName:   "テストユーザー",
		},
		{
			name:           "トークンなし",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "不正なトークン",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "期限切れトークン",
			token:          generateExpiredToken(t, "user123", "テストユーザー"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のハンドラー関数を作成
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedStatus == http.StatusOK {
					// コンテキストからユーザーIDを取得
					userID, ok := GetUserID(r.Context())
					if !ok {
						t.Error("Expected userID in context, but not found")
					} else if userID != tt.expectedUserID {
						t.Errorf("Expected userID %s, got %s", tt.expectedUserID, userID)
					}

					// コンテキストからユーザー名を取得
					userName, ok := GetUserName(r.Context())
					if !ok {
						t.Error("Expected userName in context, but not found")
					} else if userName != tt.expectedName {
						t.Errorf("Expected userName %s, got %s", tt.expectedName, userName)
					}
				}
				w.WriteHeader(http.StatusOK)
			})

			// ミドルウェアをテストハンドラーに適用
			handler := AuthMiddleware(testHandler)

			// リクエストを作成
			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			// レスポンスレコーダーを作成
			w := httptest.NewRecorder()

			// ハンドラーを実行
			handler.ServeHTTP(w, req)

			// レスポンスを検証
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// 有効なトークンを生成するヘルパー関数
func generateValidToken(t *testing.T, userID string, userName string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   userID,
		"username": userName,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	return tokenString
}

// 期限切れトークンを生成するヘルパー関数
func generateExpiredToken(t *testing.T, userID string, userName string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   userID,
		"username": userName,
		"exp":      time.Now().Add(-time.Hour).Unix(), // 過去の時間を指定
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	return tokenString
}

// コンテキスト関数のテスト
func TestGetUserFunctions(t *testing.T) {
	// テスト用のコンテキストを作成
	testUserID := "test-user-id"
	testUserName := "test-user-name"

	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, testUserID)
	ctx = context.WithValue(ctx, userNameKey, testUserName)

	// GetUserID のテスト
	t.Run("GetUserID", func(t *testing.T) {
		userID, ok := GetUserID(ctx)
		if !ok {
			t.Error("Expected ok to be true, got false")
		}
		if userID != testUserID {
			t.Errorf("Expected userID to be %s, got %s", testUserID, userID)
		}

		// 値がない場合のテスト
		emptyCtx := context.Background()
		_, ok = GetUserID(emptyCtx)
		if ok {
			t.Error("Expected ok to be false for empty context, got true")
		}
	})

	// GetUserName のテスト
	t.Run("GetUserName", func(t *testing.T) {
		userName, ok := GetUserName(ctx)
		if !ok {
			t.Error("Expected ok to be true, got false")
		}
		if userName != testUserName {
			t.Errorf("Expected userName to be %s, got %s", testUserName, userName)
		}

		// 値がない場合のテスト
		emptyCtx := context.Background()
		_, ok = GetUserName(emptyCtx)
		if ok {
			t.Error("Expected ok to be false for empty context, got true")
		}
	})
}
