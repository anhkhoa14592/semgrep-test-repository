package vn.tiki.retail.data_index.presentation.report;

import io.swagger.v3.oas.annotations.Hidden;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;
import vn.tiki.retail.data_index.appcore.domain.mapper.CompetitorCrawlingReportResultMapper;
import vn.tiki.retail.data_index.appcore.domain.model.db.enumerations.OverviewReportType;
import vn.tiki.retail.data_index.appcore.service.report.CompetitorCrawlingReportResultService;
import vn.tiki.retail.data_index.appcore.service.report.CompetitorMetricOverviewStatisticsService;
import vn.tiki.retail.data_index.appcore.service.report.dto.CompetitorCrawlingReportDto;
import vn.tiki.retail.data_index.appcore.service.report.dto.CompetitorCrawlingReportResultUpdateDTO;
import vn.tiki.retail.data_index.appcore.service.report.dto.CompetitorMetricOverviewStatisticsUpdateDTO;
import vn.tiki.retail.data_index.appcore.service.report.dto.ResponseDto;
import vn.tiki.retail.data_index.appcore.service.themis.ThemisService;
import vn.tiki.retail.data_index.appcore.service.themis.dto.ThemisPermission;

import java.time.LocalDate;

@CrossOrigin(origins = "*")
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/system-report/competitor-crawling-report")
@Tag(name = "Report Competitor Crawling", description = "Manage competitor crawling reports")
@SecurityRequirement(name = "Authorization")
public class CompetitorCrawlingReportResultController {

    private final CompetitorCrawlingReportResultService competitorCrawlingReportResultService;
    private final CompetitorMetricOverviewStatisticsService competitorMetricOverviewStatisticsService;

    private final ThemisService themisService;

    private final ThemisPermission VIEW_PERMISSION = new ThemisPermission("RetailVerification:View", "trn:tiki:RetailVerification:report");
    private final ThemisPermission UPDATE_PERMISSION = new ThemisPermission("RetailVerification:Update", "trn:tiki:RetailVerification:report");

    

    @GetMapping("/overview/by-report-date")
    @Operation(
            summary = "Get overview report URL by report date",
            description = "Returns the downloadable report URL for overview statistics of a given date"
    )
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "URL fetched successfully",
                    content = @Content(schema = @Schema(implementation = ResponseDto.class))),
            @ApiResponse(responseCode = "403", description = "Forbidden"),
            @ApiResponse(responseCode = "404", description = "Overview report not found")
    })
    public Mono<ResponseDto<String>> findOverviewReportByDate(
            @RequestHeader(value = "Authorization") String userToken,
            @Parameter(description = "Report date", required = true, example = "2024-08-06")
            @RequestParam("date") @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate date
    ) {
        return this.themisService.checkPermission(userToken, VIEW_PERMISSION)
                .then(this.competitorCrawlingReportResultService.findOverviewReportByDate(date))
                .map(reportUrl -> new ResponseDto<>(true, "", reportUrl));
    }

}