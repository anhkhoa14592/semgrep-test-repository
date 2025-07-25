package X.presentation.search;


import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.ArraySchema;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.ExampleObject;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;
import reactor.core.publisher.Mono;
import X.appcore.domain.model.es.verification.VerificationDashboardIndex;
import X.appcore.service.aggregation.RetailStreamProcessorService;
import X.appcore.service.dataIndex.VerificationIndexService;
import X.appcore.service.themis.ThemisService;
import X.appcore.service.themis.dto.ThemisAuthRequestDTO;
import X.appcore.service.themis.dto.ThemisPermission;
import X.presentation.search.dto.verification.VerificationIndexRequest;
import X.presentation.search.dto.verification.VerificationIndexResponse;
import X.presentation.search.dto.verification.VerificationIndexSearchRequest;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/products")
@Tag(name = "Product Verification Index API", description = "APIs for managing product verification index")
@SecurityRequirement(name = "X-Authorization")
@AllArgsConstructor
@Slf4j
public class VerificationController {

    private final ThemisService themisService;
    private final VerificationIndexService verificationIndexService;
    private final RetailStreamProcessorService retailStreamProcessorService;

    private final ThemisPermission LISTING_PERMISSION = new ThemisPermission("RetailVerification:List", "trn:X:pricing");
    private final ThemisPermission VIEW_PERMISSION = new ThemisPermission("RetailVerification:View", "trn:X:pricing");
    private final ThemisPermission UPDATE_PERMISSION = new ThemisPermission("RetailVerification:Update", "trn:X:pricing");
    private final ThemisPermission CREATE_PERMISSION = new ThemisPermission("RetailVerification:Create", "trn:X:pricing");
    private final ThemisPermission DELETE_PERMISSION = new ThemisPermission("RetailVerification:Delete", "trn:X:pricing");

    @Operation(
            summary = "Search verification products",
            description = "Search product verification index with filters and category permission",
            requestBody = @io.swagger.v3.oas.annotations.parameters.RequestBody(required = true,
                    content = @Content(
                            mediaType = "application/json",
                            examples = @ExampleObject(
                                    name = "Default Search Example",
                                    summary = "A sample search body",
                                    value = "{\"brand_ids\":[],\"competitors\":[],\"product_type\":[2],\"seller_availabilities\":[],\"product_variant_options\":[],\"is_finished_verifying\":null,\"pageview_band\":[\"B\",\"A\",\"C\",\"D\"],\"resolved_date_range\":null,\"is_purchasable\":null,\"is_bpg\":null,\"pageview_l30d\":{},\"listBrandSave\":[],\"categories\":[],\"variant\":false,\"variant_price\":false,\"product_name_keyword\":null,\"sort\":[{\"field\":\"page_view\",\"order\":\"des\"}],\"offset\":0,\"limit\":10}"
                            )
                    ))
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Success",
                    content = @Content(array = @ArraySchema(schema = @Schema(implementation = VerificationIndexResponse.class)))),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @PostMapping("/search")
    public Mono<List<VerificationIndexResponse>> search(
        @RequestHeader(value = "X-Authorization") String userToken,
        @RequestBody VerificationIndexSearchRequest request
    ) {
        return themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(LISTING_PERMISSION.getAction())
                                .resource(LISTING_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                    if (!authResponse.getAllowed()) {
                        return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                    }
                    return verificationIndexService.search(request).collectList();
                });

    }

    @Operation(
            summary = "Count verification products",
            description = "Count verification index results with filters and category permission",
            requestBody = @io.swagger.v3.oas.annotations.parameters.RequestBody(required = true,
                    content = @Content(
                            mediaType = "application/json",
                            examples = @ExampleObject(
                                    name = "Default Search Example",
                                    summary = "A sample search body",
                                    value = "{\"brand_ids\":[],\"competitors\":[],\"product_type\":[2],\"seller_availabilities\":[],\"product_variant_options\":[],\"is_finished_verifying\":null,\"pageview_band\":[\"B\",\"A\",\"C\",\"D\"],\"resolved_date_range\":null,\"is_purchasable\":null,\"is_bpg\":null,\"pageview_l30d\":{},\"listBrandSave\":[],\"categories\":[],\"variant\":false,\"variant_price\":false,\"product_name_keyword\":null,\"sort\":[{\"field\":\"page_view\",\"order\":\"des\"}],\"offset\":0,\"limit\":10}"
                            )
                    ))
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Success", content = @Content(schema = @Schema(example = "{\"total_count\": 1234}"))),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @PostMapping("/count")
    public Mono<Map<String, Long>> count(
        @RequestHeader(value = "X-Authorization") String userToken,
        @RequestBody VerificationIndexSearchRequest request
    ) {
        return themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(LISTING_PERMISSION.getAction())
                                .resource(LISTING_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                    if (!authResponse.getAllowed()) {
                        return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                    }
                    return verificationIndexService.count(request)
                            .map(count -> Map.of("total_count", count));
                });
    }

    @Operation(
            summary = "Create verification index",
            description = "Index a new product for verification",
            requestBody = @io.swagger.v3.oas.annotations.parameters.RequestBody(required = true,
                    content = @Content(
                            mediaType = "application/json",
                            examples = @ExampleObject(
                                    name = "Default Example",
                                    summary = "A sample body",
                                    value = "{\"max_competitiveness_percent\":null,\"category_name\":\"Tiểu Thuyết Phương Đông\",\"competitor_summary\":null,\"min_competitiveness_percent\":null,\"main_image\":null,\"variant_type\":null,\"links_max_price\":null,\"site_max_price_list\":null,\"subcategory_name\":\"Book & Office Supplies\",\"created_at\":\"2025-07-01 12:02:17\",\"is_1p_available\":true,\"category_id\":67996,\"is_purchasable\":false,\"last_resolved_date\":null,\"price\":180000.0,\"pageview_l7d\":null,\"max_competitiveness_price\":null,\"id\":\"278226963\",\"category_ids\":[8322,316,839,844,67996],\"super_id\":null,\"total_verified_link\":null,\"seller_id\":1,\"master_pageview\":null,\"is_finished_verifying\":null,\"is_3p_available\":false,\"competitor_availability\":null,\"pageview_band_week\":null,\"min_competitiveness_price\":null,\"links_min_price\":null,\"product_name\":\"Con Chim Joong Bay Từ A Đến Z\",\"site_min_price_list\":null,\"brand_id\":null,\"pageview_band\":null,\"last_updated_at\":\"2025-07-01 12:02:17\",\"last_processed_at\":null,\"ssku_seller\":null,\"product_skus\":[\"9630162296261\",\"9630162296261\"],\"master_product_sku\":\"9630162296261\",\"pageview_l30d\":null,\"external_site_ids\":null}"
                            )
                    ))
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Created successfully"),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @PostMapping("/index")
    public Mono<Void> createIndex(
        @RequestHeader(value = "X-Authorization") String userToken,
        @RequestBody VerificationIndexRequest index
    ) {
        return themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(CREATE_PERMISSION.getAction())
                                .resource(CREATE_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                    if (!authResponse.getAllowed()) {
                        return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                    }
                    return verificationIndexService.createIndex(index).then();}
                );
    }

    @Operation(
            summary = "Bulk create verification indices",
            description = "Create verification indices in bulk",
            requestBody = @io.swagger.v3.oas.annotations.parameters.RequestBody(required = true,
                    content = @Content(
                            mediaType = "application/json",
                            examples = @ExampleObject(
                                    name = "Default Example",
                                    summary = "A sample body",
                                    value = "[{\"max_competitiveness_percent\":null,\"category_name\":\"Tiểu Thuyết Phương Đông\",\"competitor_summary\":null,\"min_competitiveness_percent\":null,\"main_image\":null,\"variant_type\":null,\"links_max_price\":null,\"site_max_price_list\":null,\"subcategory_name\":\"Book & Office Supplies\",\"created_at\":\"2025-07-01 12:02:17\",\"is_1p_available\":true,\"category_id\":67996,\"is_purchasable\":false,\"last_resolved_date\":null,\"price\":180000.0,\"pageview_l7d\":null,\"max_competitiveness_price\":null,\"id\":\"278226963\",\"category_ids\":[8322,316,839,844,67996],\"super_id\":null,\"total_verified_link\":null,\"seller_id\":1,\"master_pageview\":null,\"is_finished_verifying\":null,\"is_3p_available\":false,\"competitor_availability\":null,\"pageview_band_week\":null,\"min_competitiveness_price\":null,\"links_min_price\":null,\"product_name\":\"Con Chim Joong Bay Từ A Đến Z\",\"site_min_price_list\":null,\"brand_id\":null,\"pageview_band\":null,\"last_updated_at\":\"2025-07-01 12:02:17\",\"last_processed_at\":null,\"ssku_seller\":null,\"product_skus\":[\"9630162296261\",\"9630162296261\"],\"master_product_sku\":\"9630162296261\",\"pageview_l30d\":null,\"external_site_ids\":null}]"
                            )
                    ))
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Bulk indexed successfully"),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @PostMapping("/index/bulk")
    public Mono<Void> createIndexBulk(
        @RequestHeader(value = "X-Authorization") String userToken,
        @RequestBody List<VerificationIndexRequest> indexItems
    ) {
        return themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(CREATE_PERMISSION.getAction())
                                .resource(CREATE_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                    if (!authResponse.getAllowed()) {
                        return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                    }
                    return verificationIndexService.createIndexBulk(indexItems).then();
                });
    }

    @Operation(
            summary = "Get verification index by ID",
            description = "Retrieve detailed verification info for a product"
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Success", content = @Content(schema = @Schema(implementation = VerificationDashboardIndex.class))),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @GetMapping("/{id}")
    public Mono<VerificationDashboardIndex> get(
        @RequestHeader(value = "X-Authorization") String userToken,
        @PathVariable(value = "id") String id
    ) {
        return themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(VIEW_PERMISSION.getAction())
                                .resource(VIEW_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                            if (!authResponse.getAllowed()) {
                                return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                            }
                            return verificationIndexService.findById(id);}
                );
    }

    @Operation(
            summary = "Delete verification index by ID",
            description = "Delete a verification index using its ID"
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Deleted successfully"),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @DeleteMapping("/index/{id}")
    public Mono<Void> deleteIndex(
            @PathVariable(value = "id") String id
    ) {
        return verificationIndexService.deleteIndex(id);
    }

    @Operation(
            summary = "Sync master verification index by ID",
            description = "Sync a master verification index using its ID"
    )
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "Created successfully"),
            @ApiResponse(responseCode = "403", description = "Forbidden")
    })
    @PostMapping("/index/{master_product_id}/sync")
    public Mono<Void> sync(
            @RequestHeader(value = "X-Authorization") String userToken,
            @PathVariable(value = "master_product_id") Long masterProductId
    ) {
        return this.themisService.authorize(
                        ThemisAuthRequestDTO.builder()
                                .token(userToken)
                                .action(CREATE_PERMISSION.getAction())
                                .resource(CREATE_PERMISSION.getResource())
                                .build()
                )
                .flatMap(authResponse -> {
                    if (!authResponse.getAllowed()) {
                        return Mono.error(new ResponseStatusException(HttpStatus.FORBIDDEN, "Access denied"));
                    }
                    log.info("Sync verification index: {}", masterProductId);
                    return this.retailStreamProcessorService.processVerificationDashboardIndex(masterProductId)
                            .then();
                });
    }

}
