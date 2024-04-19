#include <cstdint>

struct PTPv2_Frame
{
    uint8_t             domainAndMessageType;
    uint8_t             versionPTP;
    uint16_t            messageLength;
    uint8_t             domainNumber;
    uint8_t             minorSdoId;
    uint16_t            flagField;
    uint64_t            correctionField;
    uint32_t            messageTypeSpecific;
    uint64_t            clockIdentity;
    uint16_t            sourcePortID;
    uint16_t            sequenceId;
    uint8_t             controlField;
    int8_t              logMessageInterval;
};

struct PTPv2_SyncFrame: PTPv2_Frame
{
    uint32_t            originTimestampSeconds1;
    uint16_t            originTimestampSeconds2;
    uint32_t            originTimestampNanoseconds;
};

int main() {
    PTPv2_SyncFrame sync_packet;

    sync_packet.domainAndMessageType = 0x10;
    sync_packet.versionPTP = 0x02; 
    sync_packet.messageLength = 44;
    sync_packet.domainNumber = 0;
    sync_packet.minorSdoId = 0;
    sync_packet.flagField = 0;
    sync_packet.correctionField = 0;
    sync_packet.messageTypeSpecific = 0;
    sync_packet.clockIdentity = 1;
    sync_packet.sourcePortID = 1;
    sync_packet.sequenceId = 1;
    sync_packet.controlField = 0;
    sync_packet.logMessageInterval = 0;
    sync_packet.originTimestampSeconds1 = 0;
    sync_packet.originTimestampSeconds2 = 0;
    sync_packet.originTimestampNanoseconds = 0;

    return 0;
}