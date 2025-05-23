package healthcheck

import (
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/healthcheck/mock"
	log "github.com/cloudtrust/common-service/v2/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuditEventsReporterChecker(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockAuditEventReporter = mock.NewAuditEventsReporterModule(mockCtrl)
	var mockTime = mock.NewTimeProvider(mockCtrl)
	mockTime.EXPECT().Now().Return(testTime).AnyTimes()

	t.Run("Success ", func(t *testing.T) {
		var auditEventReporterChecker = newAuditEventsReporterChecker("alias", mockAuditEventReporter, 10*time.Second, 10*time.Second, log.NewNopLogger(), mockTime)
		internalChecker := auditEventReporterChecker.(*auditEventsReporterChecker)
		mockAuditEventReporter.EXPECT().ReportEvent(gomock.Any(), gomock.Any()).Times(1)

		internalChecker.updateStatus()
		var res = internalChecker.response
		assert.NotNil(t, res.Connection)
		assert.Equal(t, "established", *res.Connection)
	})

	t.Run("Failure ", func(t *testing.T) {
		var auditEventReporterChecker = newAuditEventsReporterChecker("alias", mockAuditEventReporter, 1*time.Second, 10*time.Second, log.NewNopLogger(), mockTime)
		internalChecker := auditEventReporterChecker.(*auditEventsReporterChecker)
		mockAuditEventReporter.EXPECT().ReportEvent(gomock.Any(), gomock.Any()).Do(func(arg0 any, arg1 any) {
			time.Sleep(2 * time.Second)
		})

		internalChecker.updateStatus()
		res := internalChecker.response
		assert.NotNil(t, res.Message)
		assert.Equal(t, "Events reporter timeout", *res.Message)
	})
}
